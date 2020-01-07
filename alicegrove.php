<?php
// Alice Grove

// This one is finished, and Jeph re-hosted it as a single sequential list of
// files, so we're no longer reliant on scraping the Tumblr blog!

// Sadly there are lots of weird things about the setup here that Jeph didn't
// account for, his pages don't even fully load all the comics!

if (!is_dir('alicegrove')) {
    mkdir('alicegrove');
}

// We're using spread operators in an array here, which requires PHP 7.4+
$jpegs = [
    35, 70, 78, 83, 84, 98,
    100, 107, 113, 124, ...range(126, 132), 134, 136, 141, 145,
    153, 159, 164, ...range(168, 183), 186, 196,
];
$jpegs = array_merge($jpegs, );
for ($i = 1; $i <= 205; $i++) {
    // 109 and 165 are unique, 137 doesn't exist :P
    if (in_array($i, [109, 165, 137])) continue;
    $ext = in_array($i, $jpegs) ? 'jpg' : 'png';
    $path = "alicegrove/$i.$ext";
    if (!is_file($path)) {
        $url = "https://www.questionablecontent.net/images/alice/$i.$ext";
        echo "Downloading $i.$ext\n";
        file_put_contents($path, file_get_contents($url));
        usleep(5e5);
    }
}

// Handle non-standard images
$extra = ['109-1.jpg', '109-2.png', '165-1.png', '165-2.jpg'];
foreach ($extra as $img) {
    $path = "alicegrove/$img";
    if (!is_file($path)) {
        $url = "https://www.questionablecontent.net/images/alice/$img";
        echo "Downloading $img\n";
        file_put_contents($path, file_get_contents($url));
        usleep(5e5);
    }
}
