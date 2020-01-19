# Squaregame
Puzzle game written in Go

So the idea is you have a grid of squares with two components:
* square type (what it does when you click)
* square state

type could be:
* toggles state of 4conn/8conn
* toggles state of whole row/col
* autonomous
  * every click
  * every neighbouring update (this could get into CA land, like you click and it will do some toggling along the line)

honestly it would be so easy to make puzzles by just randomly clicking and then they have to unscramble. Its like doing RSA. Hard to separate.

as for states I was thinking just 0,1


UI: grid lines, inner square colour is state, outer square colour is type, gently pulsating selection square, number of levels

