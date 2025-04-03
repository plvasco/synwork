🧱 Using AWS CodeArtifact
AWS CodeArtifact is a fully managed artifact repository service that makes it easy to securely store, publish, and share software packages used in your development process.

🔐 Prerequisites
Before using CodeArtifact, ensure you have:

An AWS account with permissions to access CodeArtifact.

AWS CLI configured (aws configure).

Your AWS region and repository domain name.

🔧 Step 1: Set Up a Domain and Repository
Create a domain:

bash
Copy
Edit
aws codeartifact create-domain --domain my-domain
Create a repository:

bash
Copy
Edit
aws codeartifact create-repository --domain my-domain --repository my-repo
🔑 Step 2: Authenticate
To authenticate, use the following command to get a temporary token (valid for 12 hours):

bash
Copy
Edit
aws codeartifact get-authorization-token --domain my-domain --query authorizationToken --output text
📦 Step 3: Using with Package Managers
npm (Node.js)
Configure npm:

bash
Copy
Edit
aws codeartifact login --tool npm --domain my-domain --repository my-repo
Use npm as usual:

bash
Copy
Edit
npm install <your-package>
pip (Python)
Configure pip:

bash
Copy
Edit
aws codeartifact login --tool pip --domain my-domain --repository my-repo
Use pip:

bash
Copy
Edit
pip install <your-package>
Maven (Java)
Configure Maven:

bash
Copy
Edit
aws codeartifact login --tool maven --domain my-domain --repository my-repo
Add to pom.xml:

xml
Copy
Edit
<repositories>
  <repository>
    <id>my-repo</id>
    <url>https://<domain>-<account_id>.d.codeartifact.<region>.amazonaws.com/maven/my-repo/</url>
  </repository>
</repositories>
📤 Step 4: Publishing Packages
Each package manager has its own method. Example for npm:

bash
Copy
Edit
npm publish
Make sure your package.json is configured with the correct registry URL.

🧼 Step 5: Clean Up (Optional)
To delete the repository and domain:

bash
Copy
Edit
aws codeartifact delete-repository --domain my-domain --repository my-repo
aws codeartifact delete-domain --domain my-domain
