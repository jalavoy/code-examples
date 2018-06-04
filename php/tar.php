// this is just some useful code I wrote to tar up directories while shipping them to a client through a webserver. 
// I use this mostly to share files with my friends without having to tar them up manually before giving them a link.
<?php
    $file = $_REQUEST['file'];

    if ( ! file_exists('/home/peter.v3trae.net/' . $file) ) {
        header("HTTP/1.1 500 Internal Server Error");
        exit();
    }
    if ( preg_match('/\.(\.)?/', $file) ) {
        header("HTTP/1.1 500 Internal Server Error");
        exit();
    }
    if ( ! preg_match('/^[a-zA-Z0-9\-\_\.]+$/', $file) ) {
        header("HTTP/1.1 500 Internal Server Error");
        exit();
    }


    header('Content-type: application/gzip');
    header('Content-Disposition: attachment; filename="' . $file . '.tar.gz"');
    $file = escapeshellarg($file);
    $cmd = "/usr/local/bin/gtar cz -C /home/peter.v3trae.net -f - $file";
    $fh = popen($cmd, 'r');
    while ( ! feof($fh) ) {
        print fread($fh, 8192);
    }
    pclose($fh);
?>
