<?php
// Questionable Content
// Simple to archive, get the current day's number and work in reverse.

// This one can be re-run several times and it'll safely get images missed in
// the last run if it was stopped before completing

$html = file_get_contents('https://www.questionablecontent.net/');
preg_match('@/comics/([0-9]+)\\.png@', $html, $matches);

if (!$matches[1]) {
    echo "No comic found on home page! :(\n";
    return;
}

if (!is_dir('qc')) {
    mkdir('qc');
}

$start = $matches[1];

for ($i = $start; $i > 0; $i--) {
    if (is_file("qc/$i.png") || is_file("qc/$i.jpg") || is_file("qc/$i.gif")) {
        continue;
    }

    echo "Downloading #$i\n";
    $url = "https://www.questionablecontent.net/comics/$i.png";
    $data = @file_get_contents($url);
    if ($data) {
        file_put_contents("qc/$i.png", $data);
    } else {
		$url = "https://www.questionablecontent.net/comics/$i.jpg";
		$data = @file_get_contents($url);
		if ($data) {
			file_put_contents("qc/$i.jpg", $data);
		} else {
			$url = "https://www.questionablecontent.net/comics/$i.gif";
			$data = @file_get_contents($url);
			if ($data) {
				file_put_contents("qc/$i.gif", $data);
			}
		}
	}

    usleep(5e5);
}
