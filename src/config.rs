use serde::{Deserialize, Serialize};
use std::path::PathBuf;

#[derive(Deserialize, Serialize)]
pub struct Config {
    #[serde(rename = "comicsdir", default = "default_comics_dir")]
    pub comics_dir: String,
}

fn default_comics_dir() -> String {
    "comics".to_string()
}

impl Default for Config {
    fn default() -> Self {
        Config {
            comics_dir: default_comics_dir(),
        }
    }
}

fn config_path(override_path: Option<&str>) -> Option<PathBuf> {
    if let Some(p) = override_path {
        return Some(PathBuf::from(p));
    }
    dirs::home_dir().map(|h| h.join(".config").join("comic-archiver.yaml"))
}

pub fn load(override_path: Option<&str>) -> Config {
    let path = match config_path(override_path) {
        Some(p) => p,
        None => return Config::default(),
    };

    // Write default config if it doesn't exist
    if !path.exists() {
        if let Some(parent) = path.parent() {
            std::fs::create_dir_all(parent).ok();
        }
        let default = Config::default();
        if let Ok(yaml) = serde_yaml::to_string(&default) {
            std::fs::write(&path, yaml).ok();
        }
        return default;
    }

    match std::fs::read_to_string(&path) {
        Ok(contents) => serde_yaml::from_str(&contents).unwrap_or_default(),
        Err(_) => Config::default(),
    }
}
