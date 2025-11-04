#!/bin/bash

# This script generates standardized commit messages that comply with CI requirements.
# It analyzes staged changes and suggests appropriate commit message formats based on 
# the project's conventions.
#
# Supported commit formats:
# - resource/alicloud_xxx: description
# - docs: description  
# - ci: description
# - New Resource: alicloud_xxx
# - New Data Source: alicloud_xxx
# - website: description

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color



# Analyze file changes and suggest commit type
analyze_changes() {
  local resource_files=()
  local test_files=()
  local doc_files=()
  local ci_files=()
  local website_files=()
  local other_files=()
  
  while IFS= read -r file; do
    case "$file" in
      alicloud/resource_*.go)
        if [[ "$file" != *"_test.go" ]]; then
          resource_files+=("$file")
        else
          test_files+=("$file")
        fi
        ;;
      alicloud/data_source_*.go)
        if [[ "$file" != *"_test.go" ]]; then
          resource_files+=("$file")
        else
          test_files+=("$file")
        fi
        ;;
      *.md|docs/*)
        doc_files+=("$file")
        ;;
      ci/*|.github/*|*pipeline*|*yml|*yaml|scripts/*)
        ci_files+=("$file")
        ;;
      website/*)
        website_files+=("$file")
        ;;
      *)
        other_files+=("$file")
        ;;
    esac
  done <<< "$STAGED_FILES"
  
  # Determine primary change type and return files
  if [[ ${#resource_files[@]} -gt 0 ]]; then
    echo "resource"
    for file in "${resource_files[@]}"; do
      echo "$file"
    done
  elif [[ ${#website_files[@]} -gt 0 ]]; then
    echo "website"
  elif [[ ${#doc_files[@]} -gt 0 ]]; then
    echo "docs"
  elif [[ ${#ci_files[@]} -gt 0 ]]; then
    echo "ci"
  else
    echo "other"
  fi
}

# Extract resource name from file path
extract_resource_name() {
  local file="$1"
  local basename=$(basename "$file" .go)
  
  if [[ "$basename" =~ ^resource_alicloud_(.+)$ ]]; then
    echo "${BASH_REMATCH[1]}"
  elif [[ "$basename" =~ ^data_source_alicloud_(.+)$ ]]; then
    echo "${BASH_REMATCH[1]}"
  else
    echo ""
  fi
}

# Check if this is a new resource/data source
is_new_resource() {
  local file="$1"
  # Only check if file is being added (A) - new files only
  local status=$(git diff --cached --name-status "$file" | cut -f1)
  if [[ "$status" == "A" ]]; then
    return 0
  fi
  
  return 1
}

# Analyze added fields in resource modifications
analyze_added_fields() {
  local resource_name="$1"
  local resource_file="alicloud/resource_alicloud_${resource_name}.go"
  
  # Check if the resource file exists in staged changes
  if ! echo "$STAGED_FILES" | grep -q "$resource_file"; then
    return
  fi
  
  # Get the diff and extract added schema fields
  local added_fields=()
  local diff_output
  diff_output=$(git diff --cached "$resource_file")
  
  # Look for added schema fields (lines starting with + and containing quoted field names)
  while IFS= read -r line; do
    # Match lines like: +		"field_name": {
    if [[ "$line" =~ ^\+[[:space:]]*\"([a-zA-Z_][a-zA-Z0-9_]*)\": ]]; then
      local field_name="${BASH_REMATCH[1]}"
      # Avoid duplicates and system fields
      if [[ ! " ${added_fields[*]} " =~ " ${field_name} " ]] && \
         [[ "$field_name" != "id" ]] && \
         [[ "$field_name" != "timeouts" ]] && \
         [[ "$field_name" != "creation_time" ]] && \
         [[ "$field_name" != "status" ]]; then
        added_fields+=("$field_name")
      fi
    fi
  done <<< "$diff_output"
  
  # Return comma-separated list of fields
  if [[ ${#added_fields[@]} -gt 0 ]]; then
    local result=""
    for i in "${!added_fields[@]}"; do
      if [[ $i -eq 0 ]]; then
        result="${added_fields[i]}"
      elif [[ $i -eq $((${#added_fields[@]} - 1)) ]] && [[ ${#added_fields[@]} -gt 1 ]]; then
        if [[ ${#added_fields[@]} -eq 2 ]]; then
          result="$result, ${added_fields[i]}"
        else
          result="$result, ${added_fields[i]}"
        fi
      else
        result="$result, ${added_fields[i]}"
      fi
    done
    echo "$result"
  fi
}

# Generate commit message suggestions
generate_suggestions() {
  local analyze_output
  analyze_output=$(analyze_changes)
  local lines=()
  while IFS= read -r line; do
    [[ -n "$line" ]] && lines+=("$line")
  done <<< "$analyze_output"
  
  local change_type="${lines[0]}"
  
  case "$change_type" in
    "resource")
      # Process resource files (starting from index 1)
      local processed_resources=()
      for ((i=1; i<${#lines[@]}; i++)); do
        local file="${lines[i]}"
        local resource_name=$(extract_resource_name "$file")
        if [[ -n "$resource_name" ]]; then
          # Avoid duplicate resource suggestions
          if [[ ! " ${processed_resources[*]} " =~ " ${resource_name} " ]]; then
            processed_resources+=("$resource_name")
            if is_new_resource "$file"; then
              if [[ "$file" =~ resource_ ]]; then
                echo "New Resource: alicloud_${resource_name}"
              else
                echo "New Data Source: alicloud_${resource_name}"
              fi
            else
              echo "resource/alicloud_${resource_name}: "
            fi
          fi
        fi
      done
      ;;
    "website")
      echo "website: "
      ;;
    "docs")
      echo "docs: "
      ;;
    "ci")
      echo "ci: "
      ;;
    *)
      echo "chore: "
      ;;
  esac
}

# Auto-generate commit message based on changes
generate_commit_message() {
  echo -e "${GREEN}=== Auto-generating Commit Message ===${NC}" >&2
  echo >&2
  
  # Get suggestions as array, preserving multi-word suggestions
  local suggestions_output
  suggestions_output=$(generate_suggestions)
  local suggestions=()
  while IFS= read -r line; do
    [[ -n "$line" ]] && suggestions+=("$line")
  done <<< "$suggestions_output"
  
  local commit_message=""
  
  if [[ ${#suggestions[@]} -gt 0 ]]; then
    echo -e "${BLUE}Detected changes and generating commit message...${NC}" >&2
    
    # Automatically select the best suggestion
    local selected="${suggestions[0]}"  # Use first (most relevant) suggestion
    
    if [[ "$selected" =~ ^(New Resource|New Data Source): ]]; then
      # For new resources, add period at the end
      commit_message="$selected."
      echo -e "${GREEN}✓ New resource detected: $selected${NC}" >&2
    elif [[ "$selected" =~ resource/alicloud_.*:[[:space:]]*$ ]]; then
      # For resource modifications, analyze added fields
      local resource_name=$(echo "$selected" | sed 's/resource\/alicloud_\([^:]*\):.*/\1/')
      local added_fields=$(analyze_added_fields "$resource_name")
      if [[ -n "$added_fields" ]]; then
        commit_message="${selected}Added the field $added_fields."
        echo -e "${GREEN}✓ Resource modification detected: alicloud_$resource_name (added fields: $added_fields)${NC}" >&2
      else
        commit_message="${selected}Refactored the resource and improve the docs."
        echo -e "${GREEN}✓ Resource modification detected: alicloud_$resource_name (no new fields detected)${NC}" >&2
      fi
    elif [[ "$selected" =~ ^(docs|ci|website|chore): ]]; then
      # For other types, add generic description
      if [[ "$selected" == "docs: " ]]; then
        commit_message="docs: Update documentation."
        echo -e "${GREEN}✓ Documentation changes detected${NC}" >&2
      elif [[ "$selected" == "ci: " ]]; then
        commit_message="ci: Improve CI configuration."
        echo -e "${GREEN}✓ CI/Scripts changes detected${NC}" >&2
      elif [[ "$selected" == "website: " ]]; then
        commit_message="website: Update website content."
        echo -e "${GREEN}✓ Website changes detected${NC}" >&2
      elif [[ "$selected" == "chore: " ]]; then
        commit_message="chore: General maintenance and improvements."
        echo -e "${GREEN}✓ General changes detected${NC}" >&2
      fi
    else
      commit_message="$selected"
    fi
  else
    # Fallback: analyze files directly for generic message
    echo -e "${YELLOW}No specific patterns detected, generating generic message...${NC}" >&2
    if echo "$STAGED_FILES" | grep -q "\.go$"; then
      commit_message="chore: Update Go source files."
    elif echo "$STAGED_FILES" | grep -q "\.md$"; then
      commit_message="docs: Update documentation files."
    elif echo "$STAGED_FILES" | grep -q "\.yml$\|\.yaml$"; then
      commit_message="ci: Update CI configuration files."
    else
      commit_message="chore: Update project files."
    fi
  fi
  
  # Show all detected suggestions for reference
  if [[ ${#suggestions[@]} -gt 1 ]]; then
    echo -e "${BLUE}Other detected changes:${NC}" >&2
    for i in "${!suggestions[@]}"; do
      if [[ $i -gt 0 ]]; then
        echo "  - ${suggestions[i]}" >&2
      fi
    done
  fi
  
  echo "$commit_message"
}

# Validate commit message format
validate_commit_message() {
  local message="$1"
  
  # Debug: show what we're validating
  #echo "DEBUG: Validating message: '$message'"
  
  # Check for common patterns
  if [[ "$message" =~ ^resource/alicloud_[a-z_]+:.*$ ]] || \
     [[ "$message" =~ ^docs:.*$ ]] || \
     [[ "$message" =~ ^ci:.*$ ]] || \
     [[ "$message" =~ ^website:.*$ ]] || \
     [[ "$message" =~ ^chore:.*$ ]] || \
     [[ "$message" =~ ^New\ Resource:\ alicloud_[a-z_]+\.?$ ]] || \
     [[ "$message" =~ ^New\ Data\ Source:\ alicloud_[a-z_]+\.?$ ]]; then
    return 0
  fi
  
  echo -e "${YELLOW}WARNING: Commit message doesn't follow standard format.${NC}"
  echo "Standard formats:"
  echo "  - resource/alicloud_xxx: description"
  echo "  - docs: description"
  echo "  - ci: description"
  echo "  - website: description"
  echo "  - New Resource: alicloud_xxx"
  echo "  - New Data Source: alicloud_xxx"
  echo
  echo -e "${YELLOW}Auto-continuing with non-standard format...${NC}"
  return 0
}

# Main execution
main() {
  # Check if we're in a git repository
  if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo -e "${RED}ERROR: Not in a git repository.${NC}"
    echo "Please run this script from the root of the terraform provider repository"
    exit 1
  fi

  # Check if there are staged changes
  if ! git diff --cached --quiet; then
    echo -e "${GREEN}Found staged changes. Analyzing...${NC}"
  else
    echo -e "${YELLOW}WARNING: No staged changes found.${NC}"
    echo "Please stage your changes with 'git add' before generating commit message"
    exit 1
  fi

  # Get list of staged files
  STAGED_FILES=$(git diff --cached --name-only)
  echo -e "${BLUE}Staged files:${NC}"
  echo "$STAGED_FILES"
  echo

  echo -e "${GREEN}Terraform Provider Alicloud - Commit Message Generator${NC}"
  echo "============================================================"
  echo
  
  # Show current branch
  local current_branch=$(git branch --show-current)
  echo -e "${BLUE}Current branch: ${current_branch}${NC}"
  echo
  
  # Generate commit message
  local commit_message
  commit_message=$(generate_commit_message)
  
  if [[ -z "$commit_message" ]]; then
    echo -e "${RED}ERROR: Empty commit message${NC}"
    exit 1
  fi
  
  echo
  echo -e "${GREEN}Generated commit message:${NC}"
  echo "  $commit_message"
  echo
  
  # Validate format
  if ! validate_commit_message "$commit_message"; then
    exit 1
  fi
  
  # Auto-proceed with commit (no user confirmation needed)
  echo -e "${BLUE}Auto-committing with generated message...${NC}"
  
  git commit -m "$commit_message"
  if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Commit created successfully!${NC}"
    
    # Show recent commits
    echo
    echo -e "${BLUE}Recent commits:${NC}"
    git --no-pager log --oneline -3
    exit 0
  else
    echo -e "${RED}❌ Commit failed!${NC}"
    exit 1
  fi
}

# Show help
show_help() {
  echo "Terraform Provider Alicloud - Commit Message Generator"
  echo "====================================================="
  echo
  echo "This script generates standardized commit messages based on your staged changes."
  echo
  echo "Usage: $0 [options]"
  echo
  echo "Options:"
  echo "  -h, --help    Show this help message"
  echo
  echo "Supported commit formats:"
  echo "  resource/alicloud_xxx: description    - For resource modifications"
  echo "  docs: description                     - For documentation changes"
  echo "  ci: description                       - For CI/CD changes"
  echo "  website: description                  - For website changes"
  echo "  New Resource: alicloud_xxx            - For new resources"
  echo "  New Data Source: alicloud_xxx         - For new data sources"
  echo
  echo "Examples:"
  echo "  resource/alicloud_ecs_instance: Add support for new instance types"
  echo "  docs: Update installation instructions"
  echo "  ci: Improve breaking change detection"
  echo "  New Resource: alicloud_wafv3_defense_rule"
}

# Parse command line arguments
case "${1:-}" in
  -h|--help)
    show_help
    exit 0
    ;;
  *)
    main "$@"
    ;;
esac
