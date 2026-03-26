use super::{
    download_file, extract_basename, is_absolute_url, last_downloaded_page,
    record_last_page, ArchiveError, Logger,
};
use super::generic::http_get_string;
use regex::Regex;
use std::time::Duration;

const DELAY: Duration = Duration::from_millis(500);

/// Like `generic`, but downloads **all** images matched on a page before
/// following the previous-page link (e.g. xkcd sometimes posts two panels).
pub fn multi_image_generic(
    start_url: &str,
    dir: &str,
    file_match: &Regex,
    file_prefix: &str,
    prev_link_match: &Regex,
    skip_existing: bool,
    logger: &Logger,
) {
    if let Err(e) = std::fs::create_dir_all(dir) {
        logger.println(&format!("Failed to create directory: {e}"));
        return;
    }

    let mut url = start_url.to_string();
    let mut last = last_downloaded_page(dir);

    'outer: loop {
        let html = match http_get_string(&url) {
            Ok(s) => s,
            Err(e) => {
                logger.println(&format!("Failed to load page: {e}"));
                return;
            }
        };

        // Find and download all comic images on this page (up to 2 matches)
        let all_caps: Vec<_> = file_match.captures_iter(&html).take(2).collect();
        for caps in &all_caps {
            let img_path = caps.get(1).unwrap().as_str();
            let basename = match extract_basename(img_path) {
                Some(b) => b,
                None => continue,
            };
            let dest_path = format!("{dir}/{basename}");
            let img_url = format!("{file_prefix}{img_path}");

            match download_file(&basename, &dest_path, &img_url, logger) {
                Ok(()) => {}
                Err(ArchiveError::FileExists) => {
                    if !skip_existing {
                        logger.println(&format!("File exists: {dest_path}"));
                        return;
                    }
                    if !last.is_empty() {
                        logger.println(&format!("Skipping to URL: {last}"));
                        url = last.clone();
                        last = String::new();
                        continue 'outer;
                    }
                }
                Err(e) => {
                    logger.println(&format!("Error: {e}"));
                    return;
                }
            }
        }

        // Find link to previous comic
        let link_caps = match prev_link_match.captures(&html) {
            Some(c) => c,
            None => {
                logger.println("No previous URL found");
                return;
            }
        };
        let prev_link = link_caps.get(1).unwrap().as_str();
        url = if is_absolute_url(prev_link) {
            prev_link.to_string()
        } else {
            format!("{start_url}{prev_link}")
        };
        record_last_page(dir, &url, logger);

        std::thread::sleep(DELAY);
    }
}
