<?php
// Octopus Pie

// This one's finished and currently (2019) doing daily re-runs. We'll start
// at the last original strip posted in 2017 and work our way back to the start.

$html = file_get_contents('http://www.octopuspie.com/2017-06-05/1023-1026-thats-it/');
preg_match('@/strippy/(.+\\.(jpg|png|gif))@', $html, $matches);

if (empty($matches[1])) {
	echo "No comic found on home page! :(\n";
	return;
}

if (!is_dir('octopuspie')) {
	mkdir('octopuspie');
}

// Download current page's comic, load previous comic webpage, repeat
while (true) {
	if (is_file('octopuspie/' . $matches[1])) {
		return;
	}

	echo "Downloading {$matches[1]}\n";
	$url = "http://www.octopuspie.com/strippy/{$matches[1]}";
	$data = @file_get_contents($url);
	if ($data) {
		file_put_contents("octopuspie/{$matches[1]}", $data);
	}

    // Find previous page link
	$regex = '@href="(http://www.octopuspie.com/[0-9a-zA-Z/-]+)" rel="prev"@';
	preg_match($regex, $html, $prevMatch);

	if (empty($prevMatch[1])) {
        // TODO: Fix issue loading pages before 2018-09-05
		echo "No previous URL found!\n";
		return;
	}

	$html = @file_get_contents($prevMatch[1]);
	if (!$html) {
		echo "Failed to load previous page!\n";
		return;
	}

	preg_match('@/strippy/(.+\\.(jpg|png|gif))@', $html, $matches);
	if (empty($matches[1])) {
		echo "No image found on page!\n";
		return;
	}

	usleep(5e5);
}
