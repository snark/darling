# 💔 Darling: A Filtering RSS Aggregator

Darling is an RSS aggregator—it can combine multiple RSS or Atom feeds into a single feed—that provides simple filtering. It thus can serve as a feed muxer and demuxer.

## Aggregation

Darling will transform an unlimited number of feeds into a single feed. For instance, to produce a unified feed of posts from Lobste.rs and Tilde News: `darling https://tilde.news/rss https://lobste.rs/rss`. All items are interleaved into a single feed, sorted by creation time.

## Word Matching

Darling supports blacklisting (based on case-insensitive, whole-word matching) and whitelisting. Whitelisting takes priority over blacklisting. The asterisk is a wildcard matcher, and multiple tokens to match may be provided.

To produce a list of Julia Evans blog posts that _aren't_ about SQL or her zine: `darling -b sql -b zine -b zines https://jvns.ca/atom.xml`.

To demux her feed into one that's zine-specific: `darling -b "*" -w zine -w zines https://jvns.ca/atom.xml`.

Be aware that producing an empty feed is a valid result!

## Time Matching

Darling also supports time-based matching, returning all items _published_ (not _updated_) after a given time. This time can be relative to now or a specific date or datetime. For relative times, a subset of the standard Unix date format tokens are used: `YmdHMS`.

* Every NetNewsWire commit found from the last twelve hours: `darling  --since 12H https://github.com/brentsimmons/NetNewsWire/commits/master.atom`
* Every Dinosaur Comics entry found from the last week: `darling --since 7d https://qwantz.com/rssfeed.php`
* Every Lambda the Ultimate entry found since the beginning of 2019: `darling http://lambda-the-ultimate.org/rss.xml --since 2019-01-01`
* Every Hill Cantons entry found since Halloween, 2018: `darling https://hillcantons.blogspot.com/feeds/posts/default?alt=rss --since 2018-10-31T00:00:00-04:00`

Note that darling does not currently (and may never) handle paging through older issues on feeds which support offsets; adjust your expectations accordingly when trying to load older content.

## Limits

You can also restrict the number of items process per feed. If you wanted to get only the most recent item from film critic Nathan Rabin, perhaps to support a widget: `https://www.nathanrabin.com/happy-place/?format=rss -n 1`. 

Word matching, time matching, and limits may all be applied within a single call: `darling -n 1 --since 3d --b cat --b dog --b ghost https://strangeco.blogspot.com/feeds/posts/default` would return a feed consisting of the last post from the Strange Company blog, but only if it was in the last three days and didn't mention a dog, a cat, or a ghost (and _especially_ not a ghost dog or cat).
