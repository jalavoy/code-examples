#!/usr/local/bin/perl
use strict;
use Digest::MD5 qw(md5_hex);
use Parallel::ForkManager;
my $pm = Parallel::ForkManager->new('250');
use List::Util qw(shuffle);

my ( $goal ) = @ARGV;
if ( ! $goal ) {
	die "[!] Must specify md5!\n";
}
my $start = time();
my @words;
open(my $DAT, '<', 'deps/wordlist.txt');
	while(<$DAT>) {
		chomp();
		push(@words, $_);
	}
close($DAT);
@words = shuffle(@words);

$pm->run_on_finish(
	sub {
		my ($pid, $exit_code, $ident, $exit_signal, $core_dump, $data_structure_reference) = @_;
		if ( defined($data_structure_reference) ) {
			my ( $word1, $word2, $md5 ) = split(/:/, $$data_structure_reference);
			my $finish = ( time - $start );
			print "[*] Finish running in $finish seconds. Result: $word1 $word2 [$md5]\n";
			exit();
		}
	}
);


my $count = 0;
foreach my $word1 ( @words ) {
	$count++;
	my $pid = $pm->start and next;
	foreach my $word2 ( @words ) {
		my $test = $word1 . $word2;
		my $md5 = md5_hex($test);
		if ( $md5 eq $goal ) {
			my $response = "$word1:$word2:$md5";
			$pm->finish(0, \$response);
		}
	}
	$pm->finish(1);
}
$pm->wait_all_children;
