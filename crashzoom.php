<?php
 // Crash Zoom Toon - Comics
// https://www.crashzoomtoon.com/comics

// This one is completely hard-coded right now as it relies on Wix's lightbox
// system, which is unnecessarily hard to work with outside of an actual browser

/*
The Dustbusters URLs were determined by getting the source image URLs in the
canvas elements and stripping the resampling parameters to get the originals:

let items = document.querySelectorAll('.gallery-item');
let images = [];
items.forEach(item => {
    let src = item.getAttribute('data-src');
    if (src) {
        images.push(src).replace(/\/v1\/.+/, '');
    }
});
console.log(JSON.stringify(images));

The other image URLs were manually taken from the DOM.
*/

if (!is_dir('crashzoom')) {
    mkdir('crashzoom');
}

$data = [
    'dustbusters' => [
        'https://static.wixstatic.com/media/582a0a_dfc953b98ef84d569c6d827bb553bc84~mv2_d_1748_2480_s_2.png',
        'https://static.wixstatic.com/media/582a0a_46721963409a4996a9470ce5d536b0df~mv2_d_3114_4354_s_4_2.png',
        'https://static.wixstatic.com/media/582a0a_4e1ca341d2974968937c9e5a6e17a801~mv2_d_3169_4409_s_4_2.png',
        'https://static.wixstatic.com/media/582a0a_10b1cb2523f444e0b2d721ab7365b41d~mv2_d_3169_4409_s_4_2.png',
        'https://static.wixstatic.com/media/582a0a_7482d16d4d564ec5997f0290f0a21e74~mv2_d_3169_4409_s_4_2.png',
        'https://static.wixstatic.com/media/582a0a_62622cf9ab1443ec9a72938018802f02~mv2_d_3169_4409_s_4_2.png',
        'https://static.wixstatic.com/media/582a0a_fc1cdcb7c861471bbf78c0f2be740f2d~mv2_d_3169_4409_s_4_2.png',
        'https://static.wixstatic.com/media/582a0a_379dc6c25b504706832ca0c256c86904~mv2_d_3169_4409_s_4_2.png',
        'https://static.wixstatic.com/media/582a0a_18161b1aaf8646a9b4e0a50a6edfc771~mv2_d_3169_4409_s_4_2.png',
        'https://static.wixstatic.com/media/582a0a_4f2b137aed38483f800b0d8937257152~mv2_d_3169_4409_s_4_2.png',
        'https://static.wixstatic.com/media/582a0a_eda0f8eeaff04484972131ec368825b7~mv2_d_3169_4409_s_4_2.png',
        'https://static.wixstatic.com/media/582a0a_1c3079a1df7d45db9a1ebe2831ea0a7d~mv2_d_3169_4409_s_4_2.png',
        'https://static.wixstatic.com/media/582a0a_f66997a24b7d4b3ba7e4def34304da3a~mv2_d_3169_4409_s_4_2.png',
        'https://static.wixstatic.com/media/582a0a_df57123d16684562ba52a86ec14901c1~mv2_d_3169_4409_s_4_2.png',
        'https://static.wixstatic.com/media/582a0a_a7cead1d12f34b519063789e82a32c42~mv2_d_3169_4409_s_4_2.png',
        'https://static.wixstatic.com/media/582a0a_f58a727a57354480a5eb9f05744c9013~mv2_d_3169_4409_s_4_2.png',
        'https://static.wixstatic.com/media/582a0a_b39796008749483084f328c510a58e83~mv2_d_3169_4409_s_4_2.png',
        'https://static.wixstatic.com/media/582a0a_dd456cbb336a4d8d8f56de33cfd6d0e8~mv2_d_3169_4409_s_4_2.png',
        'https://static.wixstatic.com/media/582a0a_1342b388552843689d65420aefe8b05c~mv2_d_3169_4409_s_4_2.png',
        'https://static.wixstatic.com/media/582a0a_1470051a0a7d473090de2b9f14fd6332~mv2_d_3169_4409_s_4_2.png',
        'https://static.wixstatic.com/media/582a0a_f703e3d87ab245c2a4c4ce23c983d19f~mv2_d_3169_4409_s_4_2.png',
        'https://static.wixstatic.com/media/582a0a_f89627d71ae44fe5ae1cacf6fe6d3c46~mv2_d_1748_2480_s_2.png',
    ],
    'registered' => 'https://static.wixstatic.com/media/582a0a_5b94aa724dd3411aa747dcc2cea0bd72.png',
    'jumper' => 'https://static.wixstatic.com/media/582a0a_6c4ddbe6b9104ed2972f7af2feb690c4.jpg',
    'noise-complaint' => 'https://static.wixstatic.com/media/582a0a_c98e51356d4e4b58bfead9f71834d6cf.png',
    'balloon' => 'https://static.wixstatic.com/media/582a0a_3881b0310c20474cb59b7e9169f7be83.png',
];

foreach ($data as $name => $comic) {
    if (is_array($comic)) {
        foreach ($comic as $i => $page) {
            $ext = substr($page, -3);
            file_put_contents("crashzoom/$name-$i.$ext", file_get_contents($page));
        }
    } else {
        $ext = substr($comic, -3);
        file_put_contents("crashzoom/$name.$ext", file_get_contents($comic));
    }
}
