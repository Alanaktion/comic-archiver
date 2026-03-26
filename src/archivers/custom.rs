use super::{
    download_file_wait, is_absolute_url, last_downloaded_page, record_last_page,
    ArchiveError, Logger,
};
use super::generic::http_get_string;
use regex::Regex;
use std::time::Duration;

const DELAY: Duration = Duration::from_millis(500);

/// AliceGrove has a semi-custom site with inconsistent jpg/png naming and a
/// handful of non-standard extra pages.
pub fn alice_grove(dir: &str, file_prefix: &str, end: i32, skip_existing: bool, logger: &Logger) {
    if let Err(e) = std::fs::create_dir_all(dir) {
        logger.println(&format!("Failed to create directory: {e}"));
        return;
    }

    let jpegs: &[i32] = &[
        35, 70, 78, 83, 84, 98, 100, 107, 113, 124,
        126, 127, 128, 129, 130, 131, 132, 134, 136,
        141, 145, 153, 159, 164, 168, 169, 170, 171,
        172, 173, 174, 175, 176, 177, 178, 179, 180,
        181, 182, 183, 186, 196,
    ];

    for i in 1..=end {
        // 109 and 165 are split into two parts; 137 doesn't exist
        if i == 109 || i == 165 || i == 137 {
            continue;
        }

        let ext = if jpegs.contains(&i) { "jpg" } else { "png" };
        let name = format!("{i}.{ext}");
        let dest_path = format!("{dir}/{name}");
        let img_url = format!("{file_prefix}{name}");

        match download_file_wait(&name, &dest_path, &img_url, DELAY, logger) {
            Ok(()) => {}
            Err(ArchiveError::FileExists) => {
                if !skip_existing {
                    logger.println(&format!("File exists: {dest_path}"));
                    return;
                }
            }
            Err(e) => {
                logger.println(&format!("Error: {e}"));
                return;
            }
        }
    }

    // Non-standard split pages
    for name in &["109-1.jpg", "109-2.png", "165-1.png", "165-2.jpg"] {
        let dest_path = format!("{dir}/{name}");
        let img_url = format!("{file_prefix}{name}");

        match download_file_wait(name, &dest_path, &img_url, DELAY, logger) {
            Ok(()) => {}
            Err(ArchiveError::FileExists) => {
                if !skip_existing {
                    logger.println(&format!("File exists: {dest_path}"));
                    return;
                }
            }
            Err(e) => {
                logger.println(&format!("Error: {e}"));
                return;
            }
        }
    }
}

/// Floraverse uses hash-based server filenames but human-readable page
/// identifiers in the URL, so we construct a meaningful destination filename
/// from the page path instead of the raw server filename.
pub fn floraverse(start_url: &str, dir: &str, skip_existing: bool, logger: &Logger) {
    if let Err(e) = std::fs::create_dir_all(dir) {
        logger.println(&format!("Failed to create directory: {e}"));
        return;
    }

    let file_match = Regex::new(r#"src="https://floraverse\.com/filestore/([^"]+\.(jpg|png|gif))"#)
        .unwrap();
    let file_prefix = "https://floraverse.com/filestore/";
    let prev_link_match =
        Regex::new(r#"href="(https://floraverse\.com/comic/[0-9a-zA-Z/_-]+)">◀ previous( by date|<)"#)
            .unwrap();
    let name_path_match =
        Regex::new(r#"page\.identifier = "https://floraverse\.com/comic/([0-9a-zA-Z/_-]+)""#)
            .unwrap();

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

        let file_caps = match file_match.captures(&html) {
            Some(c) => c,
            None => {
                logger.println("No comic image found");
                return;
            }
        };
        let img_server_path = file_caps.get(1).unwrap().as_str();
        let ext = file_caps.get(2).unwrap().as_str();
        let img_url = format!("{file_prefix}{img_server_path}");

        // Build a descriptive filename from the page identifier path
        let dest_name = match name_path_match.captures(&html).and_then(|c| c.get(1)) {
            Some(m) => {
                let page_path = m.as_str().trim_matches('/').replace('/', "_");
                format!("{page_path}.{ext}")
            }
            None => {
                // Fall back to server basename
                img_server_path.to_string()
            }
        };
        let dest_path = format!("{dir}/{dest_name}");

        match download_file_wait(&dest_name, &dest_path, &img_url, DELAY, logger) {
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
