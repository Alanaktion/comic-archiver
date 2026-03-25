use axum::{
    Router,
    extract::Path,
    http::{StatusCode, header},
    response::{Html, IntoResponse, Response},
    routing::get,
};

use crate::comics::COMICS;

pub async fn start(port: u16) {
    let app = Router::new()
        .route("/", get(index_handler))
        .route("/app.css", get(css_handler))
        .route("/comic/{name}", get(comic_handler))
        .route("/comic/{name}/{index}", get(comic_page_handler))
        .route("/file/{name}/{filename}", get(file_handler));

    let addr = format!("0.0.0.0:{port}");
    let listener = tokio::net::TcpListener::bind(&addr).await.expect("Failed to bind");
    eprintln!("Listening on http://localhost:{port}/");
    axum::serve(listener, app).await.expect("Server error");
}

// ──────────────────────────────────────────────────────────────────────────────
// CSS handler
// ──────────────────────────────────────────────────────────────────────────────

const APP_CSS: &str = include_str!("app.min.css");

async fn css_handler() -> impl IntoResponse {
    ([(header::CONTENT_TYPE, "text/css")], APP_CSS)
}

// ──────────────────────────────────────────────────────────────────────────────
// Index handler  –  lists all locally-available comics
// ──────────────────────────────────────────────────────────────────────────────

async fn index_handler() -> Html<String> {
    let mut comics: Vec<&str> = COMICS
        .keys()
        .filter(|k| std::path::Path::new(k).is_dir())
        .copied()
        .collect();
    comics.sort_unstable();

    let items: String = comics
        .iter()
        .map(|c| {
            format!(
                r#"<li><a class="text-blue-600 underline cursor-pointer" href="/comic/{c}">{c}</a></li>"#
            )
        })
        .collect::<Vec<_>>()
        .join("\n            ");

    Html(format!(
        r#"<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Comic Archive</title>
    <link rel="stylesheet" href="/app.css">
</head>
<body>
    <div class="container mx-auto p-4">
        <h1 class="text-2xl mb-2">Comic Archive</h1>
        <ul>
            {items}
        </ul>
    </div>
</body>
</html>"#
    ))
}

// ──────────────────────────────────────────────────────────────────────────────
// Comic handler  –  lists pages in a comic, or shows a single page
// ──────────────────────────────────────────────────────────────────────────────

async fn comic_handler(
    Path(name): Path<String>,
) -> Result<Html<String>, StatusCode> {
    render_comic_listing(&name)
}

async fn comic_page_handler(
    Path((name, index_str)): Path<(String, String)>,
) -> Result<Html<String>, StatusCode> {
    let index: usize = index_str.parse().map_err(|_| StatusCode::NOT_FOUND)?;
    render_comic_page(&name, index)
}

fn load_pages(name: &str) -> Result<Vec<String>, StatusCode> {
    if !COMICS.contains_key(name) {
        return Err(StatusCode::NOT_FOUND);
    }
    let path = std::path::Path::new(name);
    if !path.is_dir() {
        return Err(StatusCode::NOT_FOUND);
    }
    let mut pages: Vec<String> = std::fs::read_dir(path)
        .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?
        .filter_map(|e| e.ok())
        .filter(|e| e.path().is_file())
        .filter_map(|e| e.file_name().into_string().ok())
        .collect();
    pages.sort_unstable();
    Ok(pages)
}

fn render_comic_listing(name: &str) -> Result<Html<String>, StatusCode> {
    let pages = load_pages(name)?;
    let start_url = COMICS.get(name).map(|c| c.start_url).unwrap_or("");

    let items: String = pages
        .iter()
        .enumerate()
        .map(|(i, p)| {
            format!(
                r#"<li><a class="text-blue-600 underline cursor-pointer" href="/comic/{name}/{i}">{p}</a></li>"#
            )
        })
        .collect::<Vec<_>>()
        .join("\n            ");

    Ok(Html(format!(
        r#"<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{name} - Comic Archive</title>
    <link rel="stylesheet" href="/app.css">
</head>
<body>
    <div class="container mx-auto px-4 pt-2 pb-4">
        <nav class="d-flex gap-1">
            <a class="text-blue-500 hover:text-blue-600 focus-visible:text-blue-600" href="/">
                <span aria-hidden="true">&larr;</span> All comics
            </a>
        </nav>
        <div class="flex items-center gap-2 my-2">
            <h1 class="text-2xl">{name}</h1>
            <a class="text-blue-600 underline cursor-pointer" href="{start_url}">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2" aria-hidden="true">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                </svg>
                <span class="sr-only">View comic website</span>
            </a>
        </div>
        <ol class="mt-4">
            {items}
        </ol>
    </div>
</body>
</html>"#
    )))
}

fn render_comic_page(name: &str, index: usize) -> Result<Html<String>, StatusCode> {
    let pages = load_pages(name)?;
    if index >= pages.len() {
        return Err(StatusCode::NOT_FOUND);
    }

    let start_url = COMICS.get(name).map(|c| c.start_url).unwrap_or("");
    let page = &pages[index];

    let prev_nav = if index > 0 {
        let prev_page = &pages[index - 1];
        let prev_idx = index - 1;
        format!(
            r#"<a class="text-blue-500 hover:text-blue-600 focus-visible:text-blue-600" href="/comic/{name}/{prev_idx}" rel="prev">
                <span aria-hidden="true">&larr;</span> {prev_page}
            </a>"#
        )
    } else {
        "<span></span>".to_string()
    };

    let next_nav = if index + 1 < pages.len() {
        let next_page = &pages[index + 1];
        let next_idx = index + 1;
        format!(
            r#"<a class="text-blue-500 hover:text-blue-600 focus-visible:text-blue-600" href="/comic/{name}/{next_idx}" rel="next">
                {next_page} <span aria-hidden="true">&rarr;</span>
            </a>"#
        )
    } else {
        "<span></span>".to_string()
    };

    Ok(Html(format!(
        r#"<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{page} - {name} - Comic Archive</title>
    <link rel="stylesheet" href="/app.css">
    <script defer>
        document.addEventListener('keyup', (e) => {{
            if (e.keyCode == 37) {{
                const a = document.querySelector('a[rel="prev"]');
                if (a) window.location.href = a.href;
            }} else if (e.keyCode == 39) {{
                const a = document.querySelector('a[rel="next"]');
                if (a) window.location.href = a.href;
            }}
        }});
    </script>
</head>
<body>
    <div class="container mx-auto px-4 pt-2 pb-4">
        <nav class="d-flex gap-1">
            <a class="text-blue-500 hover:text-blue-600 focus-visible:text-blue-600" href="/">
                <span aria-hidden="true">&larr;</span> All comics
            </a>
            <span class="text-gray-400">/</span>
            <a class="text-blue-500 hover:text-blue-600 focus-visible:text-blue-600" href="/comic/{name}">{name}</a>
        </nav>
        <div class="flex items-center gap-2 my-2">
            <h1 class="text-2xl">{name}</h1>
            <a class="text-blue-600 underline cursor-pointer" href="{start_url}">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2" aria-hidden="true">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                </svg>
                <span class="sr-only">View comic website</span>
            </a>
        </div>
        <div class="my-4">
            <img src="/file/{name}/{page}" alt="{page}">
        </div>
        <nav class="flex justify-between">
            {prev_nav}
            {next_nav}
        </nav>
    </div>
</body>
</html>"#
    )))
}

// ──────────────────────────────────────────────────────────────────────────────
// File handler  –  serves individual image files
// ──────────────────────────────────────────────────────────────────────────────

async fn file_handler(
    Path((name, filename)): Path<(String, String)>,
) -> Result<Response, StatusCode> {
    // Validate path components to prevent directory traversal.
    // Allow only normal (non-special) path components in both segments.
    let is_safe = |s: &str| {
        let p = std::path::Path::new(s);
        p.components().all(|c| matches!(c, std::path::Component::Normal(_)))
    };
    if !is_safe(&name) || !is_safe(&filename) {
        return Err(StatusCode::BAD_REQUEST);
    }

    let path = format!("{name}/{filename}");
    let bytes = std::fs::read(&path).map_err(|_| StatusCode::NOT_FOUND)?;

    let mime = detect_mime(&bytes, &filename);
    Ok(([(header::CONTENT_TYPE, mime)], bytes).into_response())
}

/// Detect MIME type from file bytes (first 512 bytes) or fallback to extension.
fn detect_mime(bytes: &[u8], filename: &str) -> &'static str {
    // Check magic bytes first
    if bytes.starts_with(b"\x89PNG") {
        return "image/png";
    }
    if bytes.starts_with(b"\xff\xd8\xff") {
        return "image/jpeg";
    }
    if bytes.starts_with(b"GIF8") {
        return "image/gif";
    }
    if bytes.starts_with(b"RIFF") && bytes.len() >= 12 && &bytes[8..12] == b"WEBP" {
        return "image/webp";
    }
    // Fallback to extension
    match filename.rsplit('.').next().unwrap_or("") {
        "png" => "image/png",
        "jpg" | "jpeg" => "image/jpeg",
        "gif" => "image/gif",
        "webp" => "image/webp",
        _ => "application/octet-stream",
    }
}
