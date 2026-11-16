CS 372 - Part 2 Submission
Language: Go (Golang)
Name: Nguyen Do, Harshvardhan Kaul
Date: 16 November 2026

--------------------------------------------------
1. LANGUAGE STUDY LOCATION
--------------------------------------------------
The Language Study is publicly readable at:

    https://docs.google.com/document/d/1FiNDaQTrdiXgR6erxJvo0VgsQs_U93I_Fs1z_hCs59o/edit?usp=sharing

    

--------------------------------------------------
2. COMMON PROGRAM
--------------------------------------------------
File: Wilson.go

Description:
Generates a random maze using Wilson's algorithm, which builds the maze
as a uniform spanning tree of a grid using loop-erased random walks. The
program displays a description of its purpose, prompts for the number of
rows and columns, generates the maze, and prints it as text.

Output format:
  +---+   horizontal wall
  |       vertical wall
  spaces  open cell or passage

How to run (Go must be installed; see https://go.dev/dl/)[Latest version: 1.27]:

    go run Wilson.go

Then enter the number of rows and the number of columns when prompted.
Both must be positive whole numbers.

Example (5 rows x 5 columns; output differs on every run):

    +---+---+---+---+---+
    |           |       |
    +   +---+   +---+   +
    |       |           |
    +   +---+---+   +---+
    |   |       |       |
    +   +   +   +   +---+
    |       |   |       |
    +   +---+   +   +---+
    |   |       |       |
    +---+---+---+---+---+

--------------------------------------------------
3. NOTES
--------------------------------------------------

- The program does not add an entrance or exit, and does not draw
  graphics; both were optional in the assignment.