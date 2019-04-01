<?php
// Two Guys an Guy
// This one looked like it might have died for a while but is updating again!
// It's built on ComicControl (http://comicctrl.com), Hiveworks' CMS.

$html = file_get_contents('http://twogag.com/');
preg_match('@/comics/(.+\\.(jpg|png|gif))@', $html, $matches);

if (empty($matches[1])) {
	echo "No comic found on home page! :(\n";
	return;
}

if (!is_dir('twogag')) {
	mkdir('twogag');
}

// Download current page's comic, load previous comic webpage, repeat
while (true) {
	if (is_file('twogag/' . $matches[1])) {
		return;
	}

	echo "Downloading {$matches[1]}\n";
	$url = "http://twogag.com/comics/{$matches[1]}";
	$data = @file_get_contents($url);
	if ($data) {
		file_put_contents("twogag/{$matches[1]}", $data);
	}

    // Find previous page link
	$regex = '@rel="prev" href="(http://twogag.com/archives/[0-9a-zA-Z-]+)"@';
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
