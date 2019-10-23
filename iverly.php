<?php
// Iverly
// http://www.iverly.com/iverly/

// This one hasn't updated in a while, but isn't "finished" yet, so we may need
// to update the max comic ID over time. We could auto-detect it by checking the
// main Iverly page, but we can do that later if it actually updates :P

if (!is_dir('iverly')) {
    mkdir('iverly');
}

for ($i = 1; $i <= 86; $i++) {
    $name = sprintf('IVE%04d.png', $i);

    if (!is_file("iverly/{$name}")) {
        $data = file_get_contents("http://www.iverly.com/iverly/comics/{$name}");
        echo "Downloading {$name}\n";
        file_put_contents("iverly/{$name}", $data);
        usleep(5e5);
    }
}
