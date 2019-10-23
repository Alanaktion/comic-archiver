<?php
// Wigu Adventures
// https://www.wigucomics.com/adventures/

if (!is_dir('wigu-adventures')) {
    mkdir('wigu-adventures');
}

for ($i = 1; $i <= 1179; $i++) {
    $name = sprintf('WADV%04d.png', $i);

    if (!is_file("wigu-adventures/{$name}")) {
        $data = file_get_contents("https://www.wigucomics.com/adventures/comics/{$name}");
        echo "Downloading {$name}\n";
        file_put_contents("wigu-adventures/{$name}", $data);
        usleep(5e5);
    }
}
