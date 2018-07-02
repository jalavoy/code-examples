// jalavoy - 07.01.2018
// this program plays a quick guessing game with the user. It picks a number between the globals defined below, and then prompts the user to guess.
import java.util.Random;
import java.util.Scanner;
import java.util.ArrayList;
import java.util.regex.*;

public class guessing_game {
    static int min = 1;
    static int max = 100;
    static Stats s = new Stats();
    static Scanner sc = new Scanner(System.in);
    static Pattern p = Pattern.compile("^(Y|y)");

    public static void main(String[] args) {
        while (true) {
            do_game(s);
            String response = query_player("Do you want to play again?");
            if (!p.matcher(response).matches()) {
                break;
            }
        }
        s.print_stats();
        System.exit(0);
    }

    private static void do_game(Stats s) {
        int guess = -1;
        int guesscount = 0;
        int target;
        target = get_random();
        System.out.printf("\nI'm thinking of a number between %d and %d\n", min, max);
        while (guess != target) {
            guesscount++;
            String response = query_player("Your guess?");
            response = response.trim();
            guess = Integer.parseInt(response);
            if (guess == target) {
                if (guesscount == 1) {
                    System.out.println("You got it right in 1 guess!");
                } else {
                    System.out.printf("You got it right in %d guesses!\n", guesscount);
                }
            } else {
                if ( guess < target ) {
                    System.out.println("Higher!");
                } else {
                    System.out.println("Lower!");
                }
            }
        }
        s.record_stats(guesscount);
    }

    private static int get_random() {
        Random rand = new Random();
        int n = rand.nextInt(max) + 1;
        return n;
    }

    private static String query_player(String query) {
        String input;
        System.out.printf("%s ", query);
        input = sc.nextLine();
        return input;
    }    

    private static class Stats {
        int games;
        int guesses;
        ArrayList<Number> scores = new ArrayList<Number>();

        private void record_stats(int guesscount) {
            this.games++;
            this.guesses += guesscount;
            this.scores.add(guesscount);
        }

        private void print_stats() {
            System.out.println("\nOverall Results:");
            System.out.println("Total Games: " + this.games);
            System.out.println("Total Guesses: " + this.guesses);
            System.out.println("Average guesses per game: " + get_average());
            System.out.println("Best game: " + get_best());
        }

        private int get_average() {
            int total = 0;
            for (int i = 0; i < this.scores.size(); i++) {
                total += this.scores.get(i).intValue();
            }
            int average = total / this.scores.size();
            return average;
        }

        private int get_best() {
            return this.scores.get(0).intValue();
        }   
    }
}