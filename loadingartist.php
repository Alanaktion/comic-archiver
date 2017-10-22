<?php
// Loading Artist
// A bit tricky, but not too bad

// Start by navigating to the latest comic
$html = file_get_contents('https://loadingartist.com/latest');
preg_match('@/uploads/([0-9]+/[0-9]+)/([0-9a-zA-Z-]+\\.[a-z]{3,4})@', $html, $matches);

if (empty($matches[2])) {
	echo "No comic found on starting page! :(\n";
	return;
}

if (!is_dir('loadingartist')) {
	mkdir('loadingartist');
}

// Download current page's comic, load previous comic webpage, repeat
while (true) {
	if (is_file('loadingartist/' . $matches[1])) {
		return;
	}

	echo "Downloading {$matches[2]}\n";
	$url = "https://loadingartist.com/wp-content/uploads/{$matches[1]}/{$matches[2]}";
	$data = @file_get_contents($url);
	if ($data) {
		file_put_contents("loadingartist/{$matches[2]}", $data);
	}

	// Find previous page link
	$regex = '@class="normal highlight prev comic-thumb" href="(https://loadingartist.com/comic/[0-9a-zA-Z-]+/?)"@';
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

	preg_match('@/uploads/([0-9]+/[0-9]+)/([0-9a-zA-Z-]+\\.[a-z]{3,4})@', $html, $matches);

	usleep(500000);
}
