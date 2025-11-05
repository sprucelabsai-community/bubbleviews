#!/bin/bash
set -euo pipefail

current_branch="$(git rev-parse --abbrev-ref HEAD)"

if [[ -n "$(git status --porcelain)" ]]; then
	echo "Found pending changes. Preparing to commit before tagging."
	read -rp "Commit message: " commit_message
	if [[ -z "${commit_message// }" ]]; then
		echo "Commit message cannot be empty when committing pending changes." >&2
		exit 1
	fi
	git add -A
	git commit -m "$commit_message"
fi

current_version="$(tr -d '[:space:]' < VERSION)"
if [[ -z "$current_version" ]]; then
	echo "VERSION file is empty." >&2
	exit 1
fi

last_tag="$(git describe --tags --abbrev=0 2>/dev/null || true)"

commit_messages=()
if [[ -n "$last_tag" ]]; then
	log_spec="${last_tag}..HEAD"
else
	log_spec="HEAD"
fi

while IFS= read -r line; do
	[[ -z "$line" ]] && continue
	commit_messages+=("$line")
done < <(git log "$log_spec" --pretty=%s)

if [[ ${#commit_messages[@]} -eq 0 ]]; then
	echo "No commits since last tag. Nothing to publish." >&2
	exit 1
fi

determine_bump() {
	local bump="patch"
	for msg in "$@"; do
		if [[ $msg == major:* ]]; then
			echo "major"
			return
		fi
		if [[ $msg == minor:* ]]; then
			bump="minor"
		fi
	done
	echo "$bump"
}

increment_version() {
	local version="$1"
	local bump="$2"

	IFS='.' read -r major minor patch <<<"$version"
	if [[ -z $major || -z $minor || -z $patch ]]; then
		echo "Invalid semver: $version" >&2
		exit 1
	fi

	case "$bump" in
	major)
		((major += 1))
		minor=0
		patch=0
		;;
	minor)
		((minor += 1))
		patch=0
		;;
	patch)
		((patch += 1))
		;;
	*)
		echo "Unknown bump type: $bump" >&2
		exit 1
		;;
	esac

	echo "${major}.${minor}.${patch}"
}

bump_type="$(determine_bump "${commit_messages[@]}")"
new_version="$(increment_version "$current_version" "$bump_type")"
echo "$new_version" > VERSION

git add VERSION
git commit -m "release: v${new_version}"

version_tag="v${new_version}"

if git rev-parse "$version_tag" >/dev/null 2>&1; then
	echo "Tag $version_tag already exists." >&2
	exit 1
fi

git tag "$version_tag"
git push origin "$current_branch"
git push origin "$version_tag"

echo "Published $version_tag"
