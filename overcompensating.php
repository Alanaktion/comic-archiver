<?php
// Overcompensating
// http://www.overcompensating.com/oc/

if (!is_dir('overcompensating')) {
    mkdir('overcompensating');
}

for ($i = 1; $i <= 1543; $i++) {
    $name = sprintf('OC%04d.png', $i);

    if (!is_file("overcompensating/{$name}")) {
        $data = file_get_contents("http://www.overcompensating.com/oc/comics/{$name}");
        echo "Downloading {$name}\n";
        file_put_contents("overcompensating/{$name}", $data);
        usleep(5e5);
    }
}
