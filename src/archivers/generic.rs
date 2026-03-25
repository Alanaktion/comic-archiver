use super::{
    download_file_wait, extract_basename, is_absolute_url, last_downloaded_page,
    record_last_page, ArchiveError, Logger,
};
use regex::Regex;
use std::time::Duration;

const DELAY: Duration = Duration::from_millis(500);

/// Generic archiver supporting ComicPress-style sites.
///
/// Starts at `start_url` and follows previous-page links backwards until
/// it reaches the beginning of the archive (no more prev links) or until
/// it encounters a file that already exists when `skip_existing` is false.
pub fn generic(
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

    logger.println(start_url);

    let mut url = start_url.to_string();
    let mut last = last_downloaded_page(dir);

    loop {
        let html = match http_get_string(&url) {
            Ok(s) => s,
            Err(e) => {
                logger.println(&format!("Failed to load page: {e}"));
                return;
            }
        };

        // Find comic image
        let caps = match file_match.captures(&html) {
            Some(c) => c,
            None => {
                logger.println("No comic image found");
                return;
            }
        };
        let img_path = caps.get(1).unwrap().as_str();
        let img_url = format!("{file_prefix}{img_path}");
        let basename = match extract_basename(img_path) {
            Some(b) => b,
            None => {
                logger.println(&format!("Could not extract basename from: {img_path}"));
                return;
            }
        };
        let dest_path = format!("{dir}/{basename}");

        // Download image
        match download_file_wait(&basename, &dest_path, &img_url, DELAY, logger) {
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
                    continue;
                }
            }
            Err(e) => {
                logger.println(&format!("Error: {e}"));
                return;
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

/// Variant of `generic` that first resolves the actual start URL from a
/// home/landing page before beginning the backwards traversal.
pub fn generic_custom_start(
    start_url: &str,
    start_match: &Regex,
    dir: &str,
    file_match: &Regex,
    file_prefix: &str,
    prev_link_match: &Regex,
    skip_existing: bool,
    logger: &Logger,
) {
    let html = match http_get_string(start_url) {
        Ok(s) => s,
        Err(e) => {
            logger.println(&format!("Failed to load page: {e}"));
            return;
        }
    };

    let caps = match start_match.captures(&html) {
        Some(c) => c,
        None => {
            logger.println("No start URL found");
            return;
        }
    };
    let real_start = caps.get(1).unwrap().as_str().to_string();

    generic(&real_start, dir, file_match, file_prefix, prev_link_match, skip_existing, logger);
}

/// Perform an HTTP GET and return the response body as a UTF-8 string.
pub(crate) fn http_get_string(url: &str) -> Result<String, ArchiveError> {
    let resp = reqwest::blocking::get(url).map_err(|e| ArchiveError::Http(e.to_string()))?;
    resp.text().map_err(|e| ArchiveError::Http(e.to_string()))
}
