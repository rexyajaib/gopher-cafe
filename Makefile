############################# Main targets #############################
generate-buf: install-buf lint-buf generate-go-buf
########################################################################

#-------------------------------------------------------------------
# Variables															#
# -------------------------------------------------------------------
# EXCLUSION            := coincommon
BUF_GEN_GO_TEMPLATE  := buf.gen.go.yaml
BUF_GO_TEMPLATE      := buf.go.yaml
COLOR_BLUE           := "\033[1;34m"
COLOR_GREEN          := "\033[1;32m"
COLOR_RED            := "\033[1;31m"

# -------------------------------------------------------------------
# Generate protobuf for *all* services								#
# -------------------------------------------------------------------

generate-go-buf:
	@echo $(COLOR_BLUE) "🔧 Begin Generating protobuf for all services..."
	@cp $(BUF_GEN_GO_TEMPLATE) buf.gen.yaml
	@buf generate
	@cp -R github.com/rexyajaib/gopher-cafe/* .
	@rm buf.gen.yaml
	@rm -r github.com
	@echo $(COLOR_GREEN) "✅ Successfully generated all protobuf files."


lint-buf:
	@echo $(COLOR_BLUE) "🔎 Running lint for all services..." $(COLOR_RED)
	@for SERVICE in $$(find . -maxdepth 1 -type d ! -name '.*' ! -name 'java' -exec basename {} \;); do \
		cp $(BUF_GO_TEMPLATE) "$$SERVICE/buf.yaml" ; \
		if ! buf lint --path "$$SERVICE" ; then \
			echo $(COLOR_RED) "❌ Lint failed for $$SERVICE." ; \
			rm -f "$$SERVICE/buf.yaml" ; \
			exit 1 ; \
		fi ; \
		rm -f "$$SERVICE/buf.yaml" ; \
	done
	@echo $(COLOR_GREEN) "✅ The code passed the linter check."

install-buf:
	@echo $(COLOR_BLUE) "🔧 Checking if buf is installed..."
	@if ! command -v buf &> /dev/null; then \
		echo $(COLOR_RED) "❌ Buf is not installed. Installing now..."; \
		make install-buf-deps; \
	else \
		echo $(COLOR_GREEN) "✅ Buf is already installed."; \
	fi