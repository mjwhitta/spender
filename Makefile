# Include gomk if it's been checked-out: git submodule update --init
-include gomk/main.mk
-include local/Makefile

ifneq ($(unameS),windows)
spellcheck:
	@codespell -f -S ".git,local,testdata"
endif
