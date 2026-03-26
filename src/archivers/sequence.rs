use super::{download_file_wait, ArchiveError, Logger};
use std::time::Duration;

const DELAY: Duration = Duration::from_millis(500);

/// Convert a printf-style format string (e.g. `"IMG%04d.png"`) to a Rust
/// `format!`-compatible string and apply it to `i`.
fn apply_seq_pattern(pattern: &str, i: i32) -> String {
    // Handle the common zero-padded decimal specifiers used in this codebase.
    // Pattern lengths: %04d=4, %03d=4, %02d=4, %d=2 characters.
    if let Some(idx) = pattern.find("%04d") {
        return format!("{}{:04}{}", &pattern[..idx], i, &pattern[idx + 4..]);
    }
    if let Some(idx) = pattern.find("%03d") {
        return format!("{}{:03}{}", &pattern[..idx], i, &pattern[idx + 4..]);
    }
    if let Some(idx) = pattern.find("%02d") {
        return format!("{}{:02}{}", &pattern[..idx], i, &pattern[idx + 4..]);
    }
    if let Some(idx) = pattern.find("%d") {
        return format!("{}{}{}", &pattern[..idx], i, &pattern[idx + 2..]);
    }
    pattern.to_string()
}

/// Download images with sequential numeric filenames between `start` and `end`
/// (inclusive) using the given URL prefix and filename pattern.
pub fn sequential(
    dir: &str,
    file_prefix: &str,
    pattern: &str,
    start: i32,
    end: i32,
    skip_existing: bool,
    logger: &Logger,
) {
    if let Err(e) = std::fs::create_dir_all(dir) {
        logger.println(&format!("Failed to create directory: {e}"));
        return;
    }

    for i in start..=end {
        let name = apply_seq_pattern(pattern, i);
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
}
