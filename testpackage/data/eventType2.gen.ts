const eventType2s = [
	"NO_DOC",
] as const;

type EventType2 = typeof eventType2s[number];
