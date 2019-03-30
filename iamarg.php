<?php
// I AM ARG!
// Based on the Dumbing of Age ripper since they have similar DOMs
// This one is nice because every comic is named by it's ISO-8601 date!

// We start with the current day, then click the Previous link until we find
// an image we've already saved before.

$html = file_get_contents('http://iamarg.com/');
preg_match('@/comics/(.+\\.(jpg|png|gif))@', $html, $matches);

if (empty($matches[1])) {
	echo "No comic found on home page! :(\n";
	exit(1);
}

if (!is_dir('iamarg')) {
	mkdir('iamarg');
}

// Download current page's comic, load previous comic webpage, repeat
while (true) {
	if (is_file('iamarg/' . $matches[1])) {
		return;
	}

	echo "Downloading {$matches[1]}\n";
	$url = "http://iamarg.com/comics/{$matches[1]}";
	$data = @file_get_contents($url);
	if ($data) {
		file_put_contents("iamarg/{$matches[1]}", $data);
	}

	// Find previous page link
	$regex = '@href="(http://iamarg.com/[0-9a-zA-Z/-]+)" class="navi navi-prev"@';
	preg_match($regex, $html, $prevMatch);

	if (!$prevMatch[1]) {
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
