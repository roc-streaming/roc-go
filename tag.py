#!/usr/bin/env python3
import argparse
import os
import re
import subprocess

def check_version_format(version):
    pattern = re.compile(r'^\d+\.\d+\.\d+$')
    return pattern.match(version) is not None

def check_tag_exists(tag):
    try:
        subprocess.check_output(['git', 'rev-parse', tag], stderr=subprocess.STDOUT)
        return True
    except subprocess.CalledProcessError:
        return False

def write_version(version):
    with open('roc/version.go', 'r') as f:
        lines = f.readlines()

    with open('roc/version.go', 'w') as f:
        for line in lines:
            if line.startswith('var bindingsVersion ='):
                line = f'var bindingsVersion = "{version}"\n'
            f.write(line)

def commit_change(version):
    subprocess.check_call(['git', 'add', 'roc/version.go'])
    subprocess.check_call(['git', 'commit', '-m', f'Release {version}'])

def create_tag(tag):
    subprocess.check_call(['git', 'tag', tag])

def push_change(remote, tag):
    subprocess.check_call(['git', 'push', remote, 'HEAD'])
    subprocess.check_call(['git', 'push', remote, tag])

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--push', required=False, help='remote to push')
    parser.add_argument('version', help='version to release')
    args = parser.parse_args()

    version = args.version
    remote = args.push
    tag = f'v{version}'

    if not check_version_format(version):
        print(f'Error: version "{version}" is not in correct format. Correct format is "x.y.z"')
        return

    if check_tag_exists(tag):
        print(f'Error: tag "{tag}" already exists')
        return

    write_version(version)
    commit_change(version)
    create_tag(tag)
    if remote:
        push_change(remote, tag)
    print(f'Successfully released {tag}')

if __name__ == "__main__":
    main()
