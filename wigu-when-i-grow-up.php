<?php
// When I Grow Up
// https://www.wigucomics.com/whenigrowup/

if (!is_dir('wigu-when-i-grow-up')) {
    mkdir('wigu-when-i-grow-up');
}

for ($i = 1; $i <= 679; $i++) {
    $name = sprintf('WIGU%04d.jpg', $i);

    if (!is_file("wigu-when-i-grow-up/{$name}")) {
        $data = file_get_contents("https://www.wigucomics.com/whenigrowup/comics/{$name}");
        echo "Downloading {$name}\n";
        file_put_contents("wigu-when-i-grow-up/{$name}", $data);
        usleep(5e5);
    }
}
