<?php
// Sufficiently Remarkable

if (!is_dir('sufficientlyremarkable')) {
  mkdir('sufficientlyremarkable');
}

$url = 'http://sufficientlyremarkable.com/';
while ($url) {
  $html = file_get_contents($url);
  preg_match("/data-current_comic_id=['\"]([0-9]+)['\"]/", $html, $idMatch);
  preg_match('/src=[\'"]([^\'"]+)[\'"] class="comic"/', $html, $matches);
  $name = basename($matches[1]);
  if (!is_file("sufficientlyremarkable/{$idMatch[1]}-{$name}")) {
    $data = file_get_contents($matches[1]);
    echo "Downloading {$idMatch[1]}-{$name}\n";
    file_put_contents("sufficientlyremarkable/{$idMatch[1]}-{$name}", $data);
  }

  preg_match('/class="comicPagination nav-prev" href="([^"]+)"/', $html, $urlMatch);
  $url = @$urlMatch[1];

  usleep(500000);
}
