UDM-Pro Memory Monitor
======================

Patrick Wagstrom &lt;patrick@wagstrom.net&gt;

April 2021

Overview
--------

This is a really simple script that was developed for one purpose, to periodically check the amount of free memory left on my UDM-Pro and restart unifi-os if it's getting out of control.

Why would I need this? Well, one or more of the apps has a memory leak, which makes it so sometimes my system performance looks like this:

[Memory goes up, until it crashes](udm-pro-memory-crash.png)

When this happens my whole network goes down for a few minutes as the UDM-Pro has to reboot. Maybe there's some switching that can happen for things that already have resolved ARP for local connections, but nothing else. Work meetings stop. And that's just not acceptable.

This "fixes" the problem by logging into the UDM Pro on a regular basis and restarting the services if memory is below a given threshold.

License
-------

MIT License
