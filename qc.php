<?php

$html = file_get_contents('http://www.questionablecontent.net/');
preg_match('@/comics/([0-9]+).png@', $html, $matches);

if (!$matches[1]) {
	echo "No comic found on home page! :(\n";
	return;
}

if (!is_dir('qc')) {
	mkdir('qc');
}

$start = $matches[1];

for ($i = $start; $i > 0; $i--) {
	if (is_file("qc/$i.png")) {
		continue;
	}

	echo "Downloading #$i\n";
	$url = "http://www.questionablecontent.net/comics/$i.png";
	$data = @file_get_contents($url);
	if ($data) {
		file_put_contents("qc/$i.png", $data);
	}

	usleep(500000);
}
