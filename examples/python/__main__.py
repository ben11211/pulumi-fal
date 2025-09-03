import pulumi
import pulumi_fal

provider = pulumi_fal.Provider("fal", fal_key=pulumi.Config("fal").require("falKey"))

# Create a fal application equivalent to the Terraform example
sana_app = pulumi_fal.App("sana_app",
    name="sana-app",
    entrypoint="fal_demos/image/sana.py",
    git=pulumi_fal.GitConfigArgs(
        url="https://github.com/fal-ai-community/fal-demos.git",
    ))

# Export the app's properties
pulumi.export("appName", sana_app.name)
pulumi.export("revisionId", sana_app.revision_id)
pulumi.export("createdAt", sana_app.created_at)
pulumi.export("updatedAt", sana_app.updated_at)