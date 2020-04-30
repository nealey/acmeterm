This is a terminal program for [acme](https://github.com/9fans/plan9port/blob/master/man/man1/acme.1),
a text editor (in the sense that Emacs is a text editor) that was created for Plan 9.

If you're interested in acme,
check out
[a tour of acme](https://research.swtch.com/acme).

Planned Features
========

* [ ] History of stuff you've typed in, which will persist across `ssh` invocations, `ipmitool sol activate`, or whatever else
* [ ] Strip VT100/ANSI/xterm escape codes
* [ ] Understand all xterm title escapes to set window title,
  so I can see what's going on when I log into our cloud machines,
  without having to edit hundreds of .bashrc files
* [ ] Pay attention to termios to go in and out of raw mode, instead of the hack in win.c
* [ ] Helper program to let you remote edit files. It'll go something like this:
  * `cat fileYouWantToEdit`
  * select the file (maybe you can hit Esc to do this!)
  * button-2 click on the program (`remote-edit` maybe?) in your tab bar
  * a new window will open with your select contents and a guess at the remote filename
  * make your edits just like normal
  * middle-click `Put` just like normal:
    the helper program will run `cat > fileYouWantToEdit`,
    paste in your window contents, and send `^D`

Stretch Goals
---------------

* Be smart about recognizing Readline droppings and strip them out.
  Mosh plays this game extremely well;
  I might be able to do an acceptable job too.
* There's no reason why acme can't be used to actually emulate a VT100,
  with screen positioning and the whole works.
  There are even VT100 emulators already written for Go.
  The one by the Microsoft Azure team looks particularly promising.
* New escape code that means "here is a pathname and file contents. Edit this."
* Google's hterm has a cool thing where you can display images inline.
  We obviously can't display images inline,
  but we could save them off to a file and launch an image viewer.


Philosophy
========

Note To Hackers
-----------------

The
[very first commit](https://github.com/nealey/acmeterm/commit/99c54c954039bb5025b876fcaa8ac90b86d021d0)
is a working terminal.
You can send `^P` and `^N` to get to Readline command history.
It's got a lot of room for improvement,
but I used it to make the very first git commit of the code.

It is 155 lines long.

If that makes you excited,
I suggest you fork this, roll back to that commit,
and start playing around with it.


What This Does
------------------

One of the neat things acme comes with is a program called
[win](https://github.com/9fans/plan9port/blob/master/man/man1/win.1),
which starts a shell in a new Acme window, and lets you interact with it.

An aspect of win that I personally find very helpful
is that it lets you work in line mode, where nothing is sent until you hit enter.
This is extremely helpful over high-latency links,
such as my home Internet connection in April 2020,
when the entire town is working from home.
I have previously worked on
[something to provide this for emacs](https://github.com/nealey/neale-ssh.el),
but, candidly, Emacs LISP is not a language I'm very comfortable with,
and I have a love/hate relationship with Emacs.

Why Does This Exist?
------------------

Win requires you to do a few things that make you sort of a weirdo.
First, there's the `nobs` command you can set your pager to,
which deals with mostly man page formatting.
Then, there's the annoying tendency for many, many programs these days
to assume everything is capable of displaying VT100 positioning
and ANSI color codes (basically, that everything is an xterm emulator).
These codes come straight through in acme, and annoy acme users.

I never did get accustomed to having to page up through my shell output in order to re-run commands.
I got familiar with Bash's `!!`, and Bourne shell's `$_`,
and I know there are some editing commands,
but it's tiresome for me.

Finally,
my Emacs comint derivative had half-working, but buggy,
support for editing remote files
over the shell (by sending `cat >filename`),
and I thought maybe I could get this working better in acme.

I tried to get these things hacked into accessory programs for win,
but it quickly became clear I was going to have to actually hack win to do what I wanted.
But,
I mean,
that's great!
win is only about 800 lines of C,
and after I understood what was going on,
it became clear I could get almost everything on my wishlist.