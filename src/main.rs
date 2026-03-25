mod archivers;
mod comics;
mod config;
mod server;

use clap::{Parser, Subcommand};
use std::io::Write;
use std::sync::{Arc, Mutex};

#[derive(Parser)]
#[command(name = "comic-archiver", about = "A tool for archiving web comics")]
struct Cli {
    /// Config file path (default: ~/.config/comic-archiver.yaml)
    #[arg(long)]
    config: Option<String>,

    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    /// Download the specified comics from their official sources
    Archive {
        /// Download all supported comics
        #[arg(short = 'a', long)]
        all: bool,

        /// Continue partial downloads, skipping already-downloaded files
        #[arg(short = 'c', long = "continue")]
        continue_: bool,

        /// Write log output to this file instead of stderr
        #[arg(short = 'l', long)]
        log: Option<String>,

        /// Comics to archive (use `list` to see available names)
        comics: Vec<String>,
    },

    /// List supported or locally-available comics
    List {
        /// Show only comics that have been downloaded locally
        #[arg(short = 'l', long)]
        local: bool,

        /// Print one comic name per line instead of comma-separated
        #[arg(short = '1', long = "one-line")]
        one_line: bool,
    },

    /// Start a web server to browse archived comics
    Serve {
        /// Port to listen on
        #[arg(short = 'p', long, default_value = "8000")]
        port: u16,
    },
}

fn main() {
    let cli = Cli::parse();

    // Load configuration and change into the comics directory
    let cfg = config::load(cli.config.as_deref());
    println!("Using comic dir: {}", cfg.comics_dir);
    std::fs::create_dir_all(&cfg.comics_dir).expect("Failed to create comics directory");
    std::env::set_current_dir(&cfg.comics_dir).expect("Failed to change to comics directory");

    match cli.command {
        Commands::Archive { all, continue_, log, comics } => {
            run_archive(all, continue_, log, comics);
        }
        Commands::List { local, one_line } => {
            run_list(local, one_line);
        }
        Commands::Serve { port } => {
            println!("Starting server at http://localhost:{port}/");
            tokio::runtime::Builder::new_multi_thread()
                .enable_all()
                .build()
                .expect("Failed to build Tokio runtime")
                .block_on(server::start(port));
        }
    }
}

fn run_archive(all: bool, skip_existing: bool, log: Option<String>, comics: Vec<String>) {
    // Build the shared log writer
    let writer: Arc<Mutex<Box<dyn Write + Send>>> = if let Some(ref log_path) = log {
        let file = std::fs::OpenOptions::new()
            .append(true)
            .create(true)
            .open(log_path)
            .expect("Failed to open log file");
        Arc::new(Mutex::new(Box::new(file)))
    } else {
        Arc::new(Mutex::new(Box::new(std::io::stderr())))
    };

    if !all && comics.is_empty() {
        let mut w = writer.lock().unwrap();
        writeln!(w, "Specify at least one comic to download, or use --all.").ok();
        writeln!(w, "Use 'comic-archiver list' to see a list of supported comics.").ok();
        drop(w);
        std::process::exit(1);
    }

    let comic_list: Vec<String> = if all {
        comics::COMICS.keys().map(|k| k.to_string()).collect()
    } else {
        comics
    };

    let mut handles = Vec::new();

    for comic_name in comic_list {
        match comics::COMICS.get(comic_name.as_str()) {
            Some(comic) => {
                let name = comic_name.clone();
                let comic = comic.clone();
                let writer = Arc::clone(&writer);
                let handle = std::thread::spawn(move || {
                    archivers::archive(&name, &comic, skip_existing, writer);
                });
                handles.push(handle);
            }
            None => {
                let mut w = writer.lock().unwrap();
                writeln!(w, "Unknown comic: {comic_name}").ok();
            }
        }
    }

    for handle in handles {
        handle.join().ok();
    }

    let mut w = writer.lock().unwrap();
    writeln!(w, "Done.").ok();
}

fn run_list(local: bool, one_line: bool) {
    let mut comic_list: Vec<&str> = if local {
        comics::COMICS
            .keys()
            .filter(|k| std::path::Path::new(k).is_dir())
            .copied()
            .collect()
    } else {
        comics::COMICS.keys().copied().collect()
    };
    comic_list.sort_unstable();

    if one_line {
        for c in comic_list {
            println!("{c}");
        }
    } else {
        println!("{}", comic_list.join(", "));
    }
}
