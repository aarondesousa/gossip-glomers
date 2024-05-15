#!/bin/bash

# Number of times to execute the command
n=20

# Command to execute
command_to_execute="../maelstrom/maelstrom test -w broadcast --bin ./maelstrom-broadcast --node-count 5 --time-limit 20 --rate 10"

# Initialize counter for matches
num_matches=0

# Execute the command 'n' times and count matches directly
for ((i = 1; i <= n; i++)); do
    output=$($command_to_execute | grep -o "Everything" | wc -l)
    ((num_matches += output))
done

# Echo the number of matches found
echo "Number of matches found for 3b: $num_matches"

# Command to execute
# command_to_execute="../maelstrom/maelstrom test -w broadcast --bin ./maelstrom-broadcast --node-count 5 --time-limit 20 --rate 10 --nemesis partition"

# # Initialize counter for matches
# num_matches=0

# # Execute the command 'n' times and count matches directly
# for ((i = 1; i <= n; i++)); do
#     output=$($command_to_execute | grep -o "Everything" | wc -l)
#     ((num_matches += output))
# done

# # Echo the number of matches found
# echo "Number of matches found for 3c: $num_matches"