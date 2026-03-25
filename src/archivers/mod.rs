pub mod generic;
pub mod multi;
pub mod sequence;
pub mod custom;

use crate::comics::Comic;
use regex::Regex;
use std::io::Write;
use std::sync::{Arc, Mutex};
use std::time::Duration;

static PROTOCOL_MATCH: std::sync::LazyLock<Regex> =
    std::sync::LazyLock::new(|| Regex::new(r"^https?:").unwrap());

static BASENAME_MATCH: std::sync::LazyLock<Regex> =
    std::sync::LazyLock::new(|| Regex::new(r"/?([^/]+\.[A-Za-z]{3,4})$").unwrap());

/// Shared, thread-safe logger that prefixes every message.
#[derive(Clone)]
pub struct Logger {
    prefix: String,
    writer: Arc<Mutex<Box<dyn Write + Send>>>,
}

impl Logger {
    pub fn new(prefix: impl Into<String>, writer: Arc<Mutex<Box<dyn Write + Send>>>) -> Self {
        Logger {
            prefix: prefix.into(),
            writer,
        }
    }

    pub fn println(&self, msg: &str) {
        let mut w = self.writer.lock().unwrap();
        writeln!(w, "{}{}", self.prefix, msg).ok();
    }
}

/// Errors that can occur during archiving.
#[derive(Debug)]
pub enum ArchiveError {
    FileExists,
    Http(String),
    Io(std::io::Error),
    #[allow(dead_code)]
    NoImageFound,
    #[allow(dead_code)]
    NoStartUrl,
    #[allow(dead_code)]
    NoPrevUrl,
}

impl std::fmt::Display for ArchiveError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            ArchiveError::FileExists => write!(f, "file exists"),
            ArchiveError::Http(e) => write!(f, "http error: {e}"),
            ArchiveError::Io(e) => write!(f, "io error: {e}"),
            ArchiveError::NoImageFound => write!(f, "no comic image found"),
            ArchiveError::NoStartUrl => write!(f, "no start url found"),
            ArchiveError::NoPrevUrl => write!(f, "no previous url found"),
        }
    }
}

/// Download a file to `path` from `url`. Returns `ArchiveError::FileExists` if the
/// file is already present without making an HTTP request.
pub fn download_file(filename: &str, path: &str, url: &str, logger: &Logger) -> Result<(), ArchiveError> {
    if std::path::Path::new(path).exists() {
        return Err(ArchiveError::FileExists);
    }

    logger.println(&format!("Downloading: {filename}"));

    let response = reqwest::blocking::get(url)
        .map_err(|e| ArchiveError::Http(e.to_string()))?;

    let bytes = response.bytes()
        .map_err(|e| ArchiveError::Http(e.to_string()))?;

    std::fs::write(path, &bytes).map_err(ArchiveError::Io)?;

    Ok(())
}

/// Download a file and sleep `delay` after a successful transfer.
pub fn download_file_wait(filename: &str, path: &str, url: &str, delay: Duration, logger: &Logger) -> Result<(), ArchiveError> {
    download_file(filename, path, url, logger)?;
    std::thread::sleep(delay);
    Ok(())
}

/// Return the last URL saved for a comic directory, or an empty string.
pub fn last_downloaded_page(dir: &str) -> String {
    let path = format!("{dir}/.last_url");
    std::fs::read_to_string(path).unwrap_or_default()
}

/// Save the current URL so archiving can be resumed later.
pub fn record_last_page(dir: &str, url: &str, logger: &Logger) {
    let path = format!("{dir}/.last_url");
    if let Err(e) = std::fs::write(&path, url) {
        logger.println(&format!("Failed to write last URL: {e}"));
    }
}

/// Extract the basename (filename with extension) from a URL path.
pub fn extract_basename(s: &str) -> Option<String> {
    BASENAME_MATCH
        .captures(s)
        .and_then(|c| c.get(1))
        .map(|m| m.as_str().to_string())
}

/// Return true if the URL is absolute (starts with http:// or https://).
pub fn is_absolute_url(url: &str) -> bool {
    PROTOCOL_MATCH.is_match(url)
}

/// Dispatch archiving for a single comic.
pub fn archive(dir: &str, comic: &Comic, skip_existing: bool, writer: Arc<Mutex<Box<dyn Write + Send>>>) {
    let prefix = format!("[{dir}] ");
    let logger = Logger::new(prefix, writer);
    logger.println("Starting archive");

    match comic.archiver {
        "Generic" => {
            generic::generic(
                comic.start_url,
                dir,
                comic.file_match.as_ref().unwrap(),
                comic.file_prefix,
                comic.prev_link_match.as_ref().unwrap(),
                skip_existing,
                &logger,
            );
        }
        "GenericCustomStart" => {
            generic::generic_custom_start(
                comic.start_url,
                comic.start_match.as_ref().unwrap(),
                dir,
                comic.file_match.as_ref().unwrap(),
                comic.file_prefix,
                comic.prev_link_match.as_ref().unwrap(),
                skip_existing,
                &logger,
            );
        }
        "MultiImageGeneric" => {
            multi::multi_image_generic(
                comic.start_url,
                dir,
                comic.file_match.as_ref().unwrap(),
                comic.file_prefix,
                comic.prev_link_match.as_ref().unwrap(),
                skip_existing,
                &logger,
            );
        }
        "Sequential" => {
            sequence::sequential(
                dir,
                comic.file_prefix,
                comic.seq_pattern,
                comic.seq_start,
                comic.seq_end,
                skip_existing,
                &logger,
            );
        }
        "AliceGrove" => {
            custom::alice_grove(dir, comic.file_prefix, comic.seq_end, skip_existing, &logger);
        }
        "Floraverse" => {
            custom::floraverse(comic.start_url, dir, skip_existing, &logger);
        }
        other => {
            logger.println(&format!("Unknown archiver: {other}"));
        }
    }
}
