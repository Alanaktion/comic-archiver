<?php
// Girls With Slingshots

// On this one, we handle things a bit differently so that we get all of the
// color strips, as well as all of the original black-and-white ones (up to
// strip 999), while not downloading any color image twice.

$html = file_get_contents('http://girlswithslingshots.com/comic/gws-chaser-1000');
preg_match('@/comics/(.+\\.(jpg|png|gif))@', $html, $matches);

if (empty($matches[1])) {
	echo "No comic found on home page! :(\n";
	return;
}

if (!is_dir('girlswithslingshots')) {
	mkdir('girlswithslingshots');
}

// Download current page's comic, load previous comic webpage, repeat
while (true) {
	if (is_file('girlswithslingshots/' . $matches[1])) {
		return;
	}

	echo "Downloading {$matches[1]}\n";
	$url = "http://girlswithslingshots.com/comics/{$matches[1]}";
	$data = @file_get_contents($url);
	if ($data) {
		file_put_contents("girlswithslingshots/{$matches[1]}", $data);
	}

    // Find previous page link
	$regex = '@rel="prev" href="(http://girlswithslingshots.com/comic/[0-9a-zA-Z-]+)"@';
	preg_match($regex, $html, $prevMatch);

	if (empty($prevMatch[1])) {
		echo "No previous URL found!\n";
		return;
	}

	$html = @file_get_contents($prevMatch[1]);
	if (!$html) {
		echo "Failed to load previous page!\n";
		return;
	}

	preg_match('@/comics/(.+\\.(jpg|png|gif))@', $html, $matches);
	if (empty($matches[1])) {
		echo "No image found on page!\n";
		return;
	}

	usleep(5e5);
}
