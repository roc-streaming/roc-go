#!/usr/bin/env python3
import argparse
import os
import re
import subprocess
import fileinput
import sys

def check_version_format(version):
    pattern = re.compile(r'^v?\d+\.\d+\.\d+$')
    return pattern.match(version) is not None

def check_tag_exists(tag):
    try:
        subprocess.check_output(['git', 'rev-parse', tag], stderr=subprocess.STDOUT)
        return True
    except subprocess.CalledProcessError:
        return False

def write_version(version):
    print(f'Updating version in roc/version.go to {version}')
    version_line = re.compile(r'\bbindingsVersion\s*=\s*".*"')
    with fileinput.FileInput('roc/version.go', inplace=True) as f:
        for line in f:
            print(version_line.sub(f'bindingsVersion = "{version}"', line), end='')

def run_command(args):
    print(f'Running command: {" ".join(args)}')
    subprocess.check_call(args)

def commit_change(version):
    run_command(['git', 'add', 'roc/version.go'])
    run_command(['git', 'commit', '-m', f'Release {version}'])

def create_tag(tag):
    run_command(['git', 'tag', tag])

def push_change(remote, tag):
    run_command(['git', 'push', remote, 'HEAD'])
    run_command(['git', 'push', remote, tag])

def main():
    os.chdir(os.path.dirname(os.path.realpath(__file__)))

    parser = argparse.ArgumentParser()
    parser.add_argument('--push', required=False, help='remote to push')
    parser.add_argument('version', help='version to release')
    args = parser.parse_args()

    version = args.version
    if not check_version_format(version):
        print(f'Error: version "{version}" is not in correct format. Correct format is "x.y.z" or "vx.y.z"', file=sys.stderr)
        sys.exit(1)

    version = version.lstrip("v")
    remote = args.push
    tag = f'v{version}'

    if check_tag_exists(tag):
        print(f'Error: tag "{tag}" already exists', file=sys.stderr)
        sys.exit(1)

    write_version(version)
    commit_change(version)
    create_tag(tag)
    if remote:
        push_change(remote, tag)
    print(f'Successfully released {tag}')

if __name__ == "__main__":
    main()
