PROTOCOL_PATH=../client/protocol
AVDLC=$(PROTOCOL_PATH)/node_modules/.bin/avdlc

types:
	@mkdir -p kbchat/types/{keybase1,gregor1,chat1,stellar1}/
	$(AVDLC) -b -l go -t -o kbchat/types/keybase1 $(PROTOCOL_PATH)/avdl/keybase1/*.avdl
	$(AVDLC) -b -l go -t -o kbchat/types/gregor1 $(PROTOCOL_PATH)/avdl/gregor1/*.avdl
	$(AVDLC) -b -l go -t -o kbchat/types/chat1 $(PROTOCOL_PATH)/avdl/chat1/*.avdl
	$(AVDLC) -b -l go -t -o kbchat/types/stellar1 $(PROTOCOL_PATH)/avdl/stellar1/*.avdl

clean:
	rm -rf kbchat/types/
