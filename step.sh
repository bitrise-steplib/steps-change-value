#!/bin/bash

# exit if a command fails
set -e
#
# Required parameters
if [ -z "${file}" ] ; then
  echo " [!] Missing required input: file"
  exit 1
fi
if [ ! -f "${file}" ] ; then
  echo " [!] File doesn't exist at specified path: ${file}"
  exit 1
fi
if [ -z "${old_value}" ] ; then
  echo " [!] No old_value specified!"
  exit 1
fi
if [ -z "${new_value}" ] ; then
  echo " [!] No new_value specified!"
  exit 1
fi


# ---------------------
# --- Configs:

echo " (i) Provided file path: ${file}"
echo " (i) Provided old_value: ${old_value}"
echo " (i) Provided new_value: ${new_value}"
echo ""

if [ "${show_file}" == "true" ]; then
    echo
    echo "------------------------------------------"
    echo "-------------OLD  FILE--------------------"
    echo "------------------------------------------"
    cat -n "${file}"
    echo
    echo "------------------------------------------"
    echo
fi

# ---------------------
# --- Main:

{
    echo " (i) Finding line(s) with old value..."

    result="$(ls -l | grep -n "$old_value" "${file}")"
    if [ ! -z "$result" ]; then
        echo " (i) Found line(s):"
        echo "$result"
    else
        echo " (e) Old value not found"
        if [[ "${notfound_exit}" == "true" ]] ; then
            exit 1
        else
            exit 0
        fi
    fi
} || {
echo "exiting..."
}

echo " (i) Replacing..."
set -x
sed -i '' -e 's%'"$old_value"'%'"$new_value"'%g' "${file}"
set +x
echo " (i) Done"


{
    echo " (i) Finding line(s) with new value..."

    resultNew="$(ls -l | grep -n "$new_value" "${file}")"
    if [[ ! -z "$resultNew" ]] ; then
        echo " (i) Found line(s):"
        echo "$resultNew"
    else
        echo " (e) New value not found."
        echo " (e) Replacing failed"
        exit 1
    fi
} || {
echo "exiting..."
}


if [[ "${show_file}" == "true" ]] ; then
    echo
    echo "------------------------------------------"
    echo "-------------NEW  FILE--------------------"
    echo "------------------------------------------"
    cat -n "${file}"
    echo
    echo "------------------------------------------"
fi