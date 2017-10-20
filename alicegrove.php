<?php
// Alice Grove

// This is another one that's finished, so I'll hard-code a few things

// Weirdly, the easy URLs for this one actually go in *reverse*, ending with 1.
// Additionally, there are sometimes multiple images on a single page, so we
// need to account for that by detecting the Tumblr "photoset" class, which
// contains an <iframe> with the images.

// Unfortunately the photoset filenames aren't always in the right order :/

if (!is_dir('alicegrove')) {
  mkdir('alicegrove');
}

$start = 220;
$base = 'http://www.alicegrove.com/page/';

for ($i = $start; $i > 0; $i--) {
  if ($i == 1) {
    $html = file_get_contents('http://www.alicegrove.com/');
  } else {
    $html = file_get_contents('http://www.alicegrove.com/page/' . $i);
  }
  if (strpos($html, 'class="photoset"')) {
    echo "Downloading photoset...\n";
    preg_match('@src="(/post/[0-9]+/photoset_iframe/alicegrovecomic/tumblr_[^"]+)"@', $html, $frames);
    $frame = file_get_contents('http://www.alicegrove.com' . $frames[1]);
    preg_match_all('@href="(http://[0-9]+\\.media\\.tumblr\\.com/[^"]+_1280\\.png)"@', $frame, $matches);
    foreach ($matches[1] as $url) {
      $data = file_get_contents($url);
      $name = basename($url);
      echo "Downloading $name\n";
      file_put_contents("alicegrove/$name", $data);
    }
  } else {
    preg_match_all('@<figure class="photo-hires-item">\s*<a href="[^"]+"><img src="([^"]+)"@', $html, $matches);
    foreach ($matches[1] as $url) {
      $data = file_get_contents($url);
      $name = basename($url);
      echo "Downloading $name\n";
      file_put_contents("alicegrove/$name", $data);
    }
  }
}
