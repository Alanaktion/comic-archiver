<?php
// It's Walky
// This one has directory indexes! Yay!

// We basically just read the directories and filter out the thumbnails
// This is actually generalized enough it could easily fit most WP sites

// If running multiple times per day, the $years and $months can be changed
// to only contain the current year/month

$base = 'http://www.itswalky.com/wp-content/uploads/';
$years = range(2012, date('Y'));
$months = [1,12];

if (!is_dir('itswalky')) {
    mkdir('itswalky');
}

foreach ($years as $y) {
    foreach ($months as $m) {
        // Skip known nonexistent months
        if ($y == 2012 && $m < 8) {
            continue;
        }

        $dir = $base . sprintf('%d/%02d/', $y, $m);
        $html = file_get_contents($dir);

        preg_match_all('@<a href="([^/"]+)">@', $html, $matches);
        foreach ($matches[1] as $file) {
            // Remove thumbnails
            if (preg_match('/-[0-9]+x[0-9]+\\.(png|gif)$/', $file)) {
                continue;
            }

            // Remove non-comics for now
            if (!preg_match('/^[0-9]{4}/', $file)) {
                continue;
            }

            // Skip already downloaded images
            if (is_file("itswalky/$file")) {
                continue;
            }

            // Download image
            echo "Downloading $file\n";
            $url = $dir . $file;
            $data = @file_get_contents($url);
            if ($data) {
                file_put_contents("itswalky/$file", $data);
            }

            usleep(500000);
        }
    }
}
