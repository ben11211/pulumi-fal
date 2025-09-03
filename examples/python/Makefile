.PHONY: setup clean

# Build provider and install plugin
setup:
	@echo "Building provider..."
	cd ../../ && make build && cd examples/python
	@echo "Removing existing plugin..."
	pulumi plugin rm resource fal -y || true
	@echo "Installing provider plugin..."
	pulumi plugin install resource fal 0.0.0 --file ../../pulumi-provider-fal

# Clean up
clean:
	pulumi plugin rm resource fal -y || true