<?php
// CommitStrip

// I chose to use the main comic page URLs instead of the comic filenames here
// since they're actually sequential that way. Makes it much easier to read.

// We'll work from start to end

$start = 'https://www.commitstrip.com/en/2012/02/22/interview/';

if (!is_dir('commitstrip')) {
	mkdir('commitstrip');
}

$url = $start;
while ($url) {
	$html = file_get_contents($url);
	preg_match('@src="https://www.commitstrip.com/wp-content/uploads/([0-9a-zA-Z/-]+\\.[a-z]{3,4})"@', $html, $matches);
	if (!empty($matches[1])) {
		$name = str_replace('/', '-', trim(substr($url, 31), '/'));
		if (!glob("commitstrip/$name*")) {
			$data = @file_get_contents('https://www.commitstrip.com/wp-content/uploads/' . $matches[1]);
			if ($data) {
				$ext = pathinfo(parse_url($matches[1])['path'], PATHINFO_EXTENSION);
				file_put_contents("commitstrip/$name.$ext", $data);
			}
		}
	}

	preg_match('@href="(https://www.commitstrip.com/20[^"]+)" rel="next"@', $html, $matches);
	$url = str_replace('.com/20', '.com/en/20', $matches[1]);

	usleep(500000);
}
