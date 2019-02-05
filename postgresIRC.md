[00:56:25] <naquad> ysch, i had a hope that it will make things faster and partitioned heavily. hardware is currently laptop, but will aws instance (medium)
[00:56:42] <andres> RhodiumToad: I wonder if one could have a branch for short [big]integers, and then do all larger numbers in a largely branchless, fully pipelinable, way. RN the code isn't pipelineable, but if we e.g. computed val / 10, val / 100, val / 1000 etc in parallel, we could probably be a lot faster for bigger numbers.
[00:57:30] <Meraki> array*agg(tag.id), array_agg(tag.name) returns {1968424375409443863,1968424375409443862} | {Science,Politics}, is there a way to combine and return {{1968424375409443863, Science}, {968424375409443862, Politics}?
[00:58:04] 	jberkus (~smuxi@c-71-237-176-63.hsd1.or.comcast.net) joined the channel
[01:00:20] 	thallada (~thallada@fsf/member/thallada) joined the channel
[01:00:54]  <RhodiumToad>	Meraki: those seem to be incompatible data types?
[01:01:27]  <Meraki>	Yeah, one is a bigint and the other is a text
[01:01:57]  <ysch>	naquad: Partitioning is not a fairy dust that makes everything faster. More often than not, it's quite the contrary (esp. if applied blindly). Anyway... in you query, why do you have this part "SELECT market_id, MAX(open_time) AS start_time, now() AS end_time FROM ohlc_1m GROUP BY market_id", i.e. what is it for and why is it appended to the other by UNION, not UNION ALL?
[01:02:07]  <RhodiumToad>	Meraki: so they can't go in the same array
[01:02:27]  <RhodiumToad>	Meraki: you could, however, use json
[01:02:35] 	xocolatl (~Vik@212.76.254.98) joined the channel
[01:02:54]  <Meraki>	How so?
[01:03:11] 	percY- (~percY@138.68.7.139) left IRC (Remote host closed the connection)
[01:03:14]  <RhodiumToad>	how would you want the result to look?
[01:03:36]  <Meraki>	{{1968424375409443863, Science}, {968424375409443862, Politics}}
[01:03:52] <RhodiumToad> no, I already said you can't mix types in a single array
[01:04:10] <naquad> ysch, this way i'm trying to find trailing gaps. e.g. if there's no data since last record till current moment
[01:04:16] <Meraki> What is the benefit if I use json_agg?
[01:04:36] <RhodiumToad> json does allow you to mix types, which arrays don't
[01:04:54] <ysch> naquad: Also, normally partitioning is implemented to help with maintenance tasks on huge tables... so, why didn't you VACUUM ANALYZE at least some of the partitions (I assume it'll be done in production, right)?
[01:05:19] thallada (~thallada@fsf/member/thallada) left IRC (Ping timeout: 246 seconds)
[01:05:19] <ysch> naquad: Again, why UNION not UNION ALL?
[01:05:49] fs2 (~fs2@pwnhofer.at) left IRC (Quit: Ping timeout (120 seconds))
[01:05:53] <Meraki> Is I want the output to be [[1968424375409443863, Science], [968424375409443862, Politics]] using json_agg, how do I go about doing it?
[01:06:01] <naquad> ysch, it will, because it is a fresh db i didn't bother much with maintenance and focused on query itself
[01:06:11] <naquad> ysch, union all was even slower :(
[01:06:15] <naquad> no idea why so
[01:06:58] <Meraki> Rhodium Toad, please show me how to get the id and name as one object using json
[01:07:02] xocolatl (~Vik@212.76.254.98) left IRC (Ping timeout: 250 seconds)
[01:07:32] <ysch> naquad: Hmm... but result is potentially different, so which one do you actually need? Also, could you show the EXPLAIN (even without ANALYZE) of the "UNION ALL" variant?
[01:07:36] fs2 (~fs2@pwnhofer.at) joined the channel
[01:07:56] <RhodiumToad> Meraki: json_agg(json_build_array(tag.id, tag.name))
[01:08:25] <RhodiumToad> Meraki: or you could make it an object instead: would you want the ids or the names as keys?
[01:08:33] thallada (~thallada@fsf/member/thallada) joined the channel
[01:08:44] <Meraki> yes please
[01:08:59] <RhodiumToad> which?
[01:09:29] foo (~foo@unaffiliated/foo) joined the channel
[01:09:33] <RhodiumToad> json_object_agg(tag.name, tag.id) would give you {"Science":1968424375409443863, "Politics":968424375409443862}
[01:09:47] <ysch> naquad: " i didn't bother much with maintenance" --- if you just loaded this data and started playing, it may mean PostgreSQL has [almost] no statistics about you tables (and generates query plans pretty much at random), as well as no visibility maps (say goodbye to index-only scans :( ). So, it's better to do it, actually.
[01:09:53] <foo> Is there much value in using UUID postgres type? I'm having an issue with flask-admin (an admin interface) in knowing how to filter on UUID. Thinking of changing to text type. Thank you
[01:10:05] aib42 (~aib@unaffiliated/aib42) joined the channel
[01:10:30] lll7 (~kristinae@unaffiliated/lll7) left IRC (Ping timeout: 250 seconds)
[01:10:31] <Myon> text will use (more than) two times the storage because hex
[01:10:45] edrocks (~ed@096-059-219-229.res.spectrum.com) joined the channel
[01:10:51] <Meraki> Wow! didn't know that this was possible. If I want say the column name as keys, is that possible?
[01:11:29] bmomjian (~bmomjian@momjian.us) left IRC (Quit: Leaving.)
[01:11:44] <RhodiumToad> Meraki: json_agg(json_build_object('id',tag.id,'name',tag.name)) would give [{"id":1968424375409443863, "name":"Science"}, ...][01:12:11] <naquad> ysch, got it. starting vacuum
[01:12:27] zaherdirkey (~zaherdirk@37.48.153.52) left IRC (Read error: Connection reset by peer)
[01:13:22] <Meraki> I'm so happy right now. Thanks Rhodium Toad, you're a lifesaver.
[01:13:27] thallada (~thallada@fsf/member/thallada) left IRC (Ping timeout: 240 seconds)
[01:13:38] yaw (~yawboakye@37.228.248.226) left IRC (Ping timeout: 272 seconds)
[01:13:45] <foo> Meraki: Now you owe RhodiumToad a life debt.
[01:13:49] fs2 (~fs2@pwnhofer.at) left IRC (Quit: Ping timeout (120 seconds))
[01:14:31] <Meraki> I honoured, he's like a wizard to someone like me
[01:14:40] <Meraki> \*I'm
[01:16:02] fs2 (~fs2@pwnhofer.at) joined the channel
[01:16:42] xocolatl (~Vik@212.76.254.98) joined the channel
[01:17:07] nikio* (~nikio\_@unaffiliated/nikio/x-5064535) left IRC (Ping timeout: 240 seconds)
[01:17:25] thallada (~thallada@fsf/member/thallada) joined the channel
[01:17:34] <Meraki> Rhodium Toad, one more question please... How do I handle the nulls, I mean, can I use coalesce here?
[01:18:22] <RhodiumToad> sure, but if you leave the nulls alone, they'll come out as json null
[01:18:49] <RhodiumToad> what specific nulls do you need to handle?
[01:19:25] <Meraki> The bigint nulls as 0 and the text nulls as ""
[01:19:31] melissa666 (~melissa66@2601:603:4d00:18ec::d6d1) left IRC (Remote host closed the connection)
[01:19:42] <RhodiumToad> ok, easy enough with coalesce
[01:19:53] <Meraki> Thanks again
[01:20:10] <RhodiumToad> json_agg(json_build_object('id',coalesce(tag.id,0),'name',coalesce(tag.name,''))) or whatever
[01:21:13] TomTom (uid45892@gateway/web/irccloud.com/x-bhceigclndayrsqj) joined the channel
[01:22:12] xocolatl (~Vik@212.76.254.98) left IRC (Ping timeout: 250 seconds)
[01:22:33] thallada (~thallada@fsf/member/thallada) left IRC (Ping timeout: 245 seconds)
[01:22:51] <Meraki> Does MySQL support these as well, or something equivalent?
[01:23:39] <RhodiumToad> no idea
[01:24:48] xocolatl (~Vik@212.76.254.98) joined the channel
[01:25:41] <Meraki> Thanks again. I learned so much.
