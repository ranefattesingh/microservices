import os
import json
import subprocess
import urllib.request
import urllib.error
from google import genai

def post_github_comment(repo, pr_number, token, body):
    url = f"https://api.github.com/repos/{repo}/issues/{pr_number}/comments"
    headers = {
        "Authorization": f"token {token}",
        "Accept": "application/vnd.github.v3+json",
        "User-Agent": "python-urllib",
        "Content-Type": "application/json"
    }
    data = json.dumps({"body": body}).encode("utf-8")
    req = urllib.request.Request(url, data=data, headers=headers, method="POST")
    try:
        with urllib.request.urlopen(req) as response:
            return response.read()
    except urllib.error.HTTPError as e:
        error_text = e.read().decode("utf-8")
        raise RuntimeError(f"GitHub API error ({e.code}): {error_text}")


def main():
    # 1. Grab environment variables provided by GitHub Actions
    repo = os.environ.get("GITHUB_REPOSITORY")
    token = os.environ.get("GITHUB_TOKEN")
    event_path = os.environ.get("GITHUB_EVENT_PATH")
    api_key = os.environ.get("GEMINI_API_KEY")

    if not event_path:
        print("No event path found. Skipping.")
        return

    # Extract the PR number from the GitHub event metadata
    with open(event_path, "r") as f:
        event_data = json.load(f)

    pr_number = event_data.get("pull_request", {}).get("number")
    if not pr_number:
        print("This action wasn't triggered by a Pull Request. Skipping.")
        return

    # 2. Get the actual code changes
    try:
        diff = subprocess.check_output(
            ["git", "diff", "origin/main...HEAD"],
            text=True
        )
    except subprocess.CalledProcessError as e:
        print(f"Failed to get git diff: {e}")
        return

    if not diff.strip():
        print("No code changes found to review.")
        return

    # 3. Prompt the AI using the official Google GenAI SDK
    client = genai.Client(api_key=api_key)

    prompt = f"""
    You are an elite software engineer reviewing a pull request.
    Analyze this git diff for critical bugs, security risks, or architectural flaws.

    Provide your response in clean Markdown.
    Use bullet points for issues, and wrap code suggestions in appropriate code blocks.
    Be concise, polite, and direct.

    ```diff
    {diff}
    ```
    """

    print("Sending diff to Gemini...")
    response = client.models.generate_content(
        model='gemini-2.0-flash',
        contents=prompt,
    )

    ai_review = response.text

    # 4. Post the review directly onto the PR
    print(f"Posting review comment to PR #{pr_number}...")
    comment_body = f"🤖 **AI Code Review Optional Feedback:**\n\n{ai_review}"
    post_github_comment(repo, pr_number, token, comment_body)
    print("Review successfully posted!")

if __name__ == "__main__":
    main()
