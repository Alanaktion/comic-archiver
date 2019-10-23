<?php
// Wigu Havin' Fun
// https://www.wigucomics.com/fun/

if (!is_dir('wigu-havin-fun')) {
    mkdir('wigu-havin-fun');
}

for ($i = 1; $i <= 61; $i++) {
    $name = sprintf('WOO%04d.png', $i);

    if (!is_file("wigu-havin-fun/{$name}")) {
        $data = file_get_contents("https://www.wigucomics.com/fun/comics/{$name}");
        echo "Downloading {$name}\n";
        file_put_contents("wigu-havin-fun/{$name}", $data);
        usleep(5e5);
    }
}
