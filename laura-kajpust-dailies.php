<?php
// Laura Kajpust/Falconer's Dailies
// This one requires some weird stuff relative to the other archivers.

// We start with the current day, then click the Previous link until we find
// an image we've already saved before.

$html = file_get_contents('http://falcdaily.smackjeeves.com/');
preg_match('@Date Posted:</strong> ([^<>]+)</div>@', $html, $dateMatches);
preg_match('@src="(https?:)?//((www|img[0-9]+).smackjeeves.com/images/uploaded/comics/[^"]+\\.(png|jpg|gif))@', $html, $matches);

if (empty($matches[2])) {
	echo "No comic found on home page! :(\n";
	return;
}
$date = date('Y-m-d-His', strtotime($dateMatches[1]));

if (!is_dir('laura-kajpust-dailies')) {
	mkdir('laura-kajpust-dailies');
}

// Download current page's comic, load previous comic webpage, repeat
while (true) {
    $ext = substr($matches[2], -3);
    $name = basename($matches[2]);
    $path = "laura-kajpust-dailies/$date-$name.$ext";
	if (is_file($path)) {
		return;
	}

	echo "Downloading {$matches[2]}\n";
	$url = "http://{$matches[2]}";
	$data = @file_get_contents($url);
	if ($data) {
		file_put_contents($path, $data);
	}

    // Find previous page link
	$regex = '@href="(http://falcdaily.smackjeeves.com/comics/[0-9a-zA-Z/-]+)"><i class="fa fa-angle-left"@';
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

    preg_match('@Date Posted:</strong> ([^<>]+)</div>@', $html, $dateMatches);
    preg_match('@src="(https?:)?//((www|img[0-9]+).smackjeeves.com/images/uploaded/comics/[^"]+\\.(png|jpg|gif))@', $html, $matches);
	if (empty($matches[2])) {
		echo "No image found on page!\n";
		return;
    }
    $date = date('Y-m-d-His', strtotime($dateMatches[1]));

	usleep(500000);
}
