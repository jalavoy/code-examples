#!/usr/local/bin/perl
use strict;
use Data::Dumper;

my ( $in ) = @ARGV;
my @words = split(/;/, $in);
my $thing;
my @solved;
foreach my $word ( @words ) {
	my $fork = 'perl -e \'$\=$/;sub c{join"",sort split//,lc pop}$l=c(pop);for(<>){chop;c(lc$_)eq$l&&print}\' < deps/dictionary.txt ' . $word;
	chomp(my $result = `$fork`);
	push(@solved, $result);
}
my $result;
foreach my $solve ( @solved ) {
	$result .= "$solve;";
}
$result =~ s/;$//g;
print $result . "\n";
