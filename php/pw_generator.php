<!DOCTYPE html>
<html>
<head>
<META HTTP-EQUIV="Expires" CONTENT="-1">
<META HTTP-EQUIV="Expires" CONTENT="Fri, May 13 1981 08:20:00 GMT">
<META HTTP-EQUIV="Pragma" CONTENT="no-cache">
<META HTTP-EQUIV="Cache-Control" CONTENT="no-cache, must-revalidate">
<font face="monospace">
<?
    $length = array(8, 16, 24, 32, 64, 128);
    $alpha = 'a b c d e f g h i j k l m n o p q r s t u v w x y z A B C D E F G H I J K L M N O P Q R S T U V';
    $numeric = '1 2 3 4 5 6 7 8 9 0';
    $special = '! @ # $ % ^ & * ( ) _ + - = ; : , . / ? ~ `';
    $extended = '€ ‚ ƒ „ … † ‡ ˆ ‰ Š ‹ Œ Ž ‘ ’ “ ” • – — ˜ ™ š › œ ž Ÿ &#160; ¡ ¢ £ ¤ ¥ ¦ § ¨ © ª « ¬ &#173; ® ¯ ° ± ² ³ ´ µ ¶ · ¸ ¹ º » ¼ ½ ¾ ¿ À Á Â Ã Ä Å Æ Ç È É Ê Ë Ì Í Î Ï Ð Ñ Ò Ó Ô Õ Ö × Ø Ù Ú Û Ü Ý Þ ß à á â ã ä å æ ç è é ê ë ì í î ï ð ñ ò ó ô õ ö ÷ ø ù ú û ü ý þ ÿ';
    $alphaarray = explode(' ', $alpha);
    $alphanumarray = array_merge(explode(' ', $alpha), explode(' ', $numeric));
    $alphanumspecialarray = array_merge(explode(' ', $alpha), explode(' ', $numeric), explode(' ', $special));
    $alphanumspecialextendedarray = array_merge(explode(' ', $alpha), explode(' ', $numeric), explode(' ', $special), explode(' ', $extended));
?>
</head>


<body>
    <table style="width:100%">
        <tr>
            <td><b>Alpha</b></td>
        </tr>
        <?
            foreach ( $length as $i ) {
                $string = get_random_string($alphaarray, $i);
                print "<tr><td>$i chars:</td><td> $string </td>\n";
            }
        ?>
        </tr>
    </table>
    <br>
    <table style="width:100%">
        <tr>
            <td><b>Alphanum</b></td>
        </tr>
        <?
            foreach ( $length as $i ) {
                $string = get_random_string($alphanumarray, $i);
                print "<tr><td>$i chars:</td><td> $string </td>\n";
            }
        ?>
        </tr>
    </table>
    <br>
    <table style="width:100%">
        <tr>
            <td><b>Special</b></td>
        </tr>
        <?
            foreach ( $length as $i ) {
                $string = get_random_string($alphanumspecialarray, $i);
                print "<tr><td>$i chars:</td><td> $string </td>\n";
            }
        ?>
        </tr>
    </table>
    <br>
    <table style="width:63%">
        <tr>
            <td><b>Extended</b></td>
        </tr>
        <?
            foreach ( $length as $i ) {
                $string = get_random_string($alphanumspecialextendedarray, $i);
                print "<tr><td>$i chars:</td><td> $string </td>\n";
            }
        ?>
    </table>
</body>






<?
    function get_random_string($chars, $length) {
        $random_string = '';
        $string = implode('', $chars);
        $num_valid = strlen($string);
        for ( $i = 0; $i < $length; $i++ ) {
            $random_pick = mt_rand(1, $num_valid);
            $random_char = $chars[$random_pick - 1];
            $random_string .= $random_char;
        }
        return($random_string);
    }
?>
