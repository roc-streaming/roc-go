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
    print(f'--- Updating version in roc/version.go to {version}')
    version_line = re.compile(r'\bbindingsVersion\s*=\s*".*"')
    with fileinput.FileInput('roc/version.go', inplace=True) as f:
        for line in f:
            print(version_line.sub(f'bindingsVersion = "{version}"', line), end='')

def run_command(args):
    print(f'--- Running command: {" ".join(args)}')
    subprocess.check_call(args)

def commit_change(version, force):
    run_command(['git', 'add', 'roc/version.go'])
    try:
        run_command(['git', 'commit', '-m', f'Release {version}'])
    except subprocess.CalledProcessError:
        if force:
            print('--- Nothing to commit, but --force was used')
        else:
            raise

def create_tag(tag, force):
    if force:
        run_command(['git', 'tag', '-f', tag])
    else:
        run_command(['git', 'tag', tag])


def push_change(remote, tag, force):
    run_command(['git', 'push', remote, 'HEAD'])
    if force:
        run_command(['git', 'push', '-f', remote, tag])
    else:
        run_command(['git', 'push', remote, tag])

def main():
    os.chdir(os.path.dirname(os.path.realpath(__file__)))

    parser = argparse.ArgumentParser(description='Create and push tag for new release')
    parser.add_argument('--push', metavar='REMOTE', required=False,
                        help='remote to push (e.g. "origin")')
    parser.add_argument('-f', '--force', action='store_true', required=False,
                        help='force creation of tag, even when already exists')
    parser.add_argument('version', metavar='VERSION',
                        help='version to release ("N.N.N" or "vN.N.N")')
    args = parser.parse_args()

    version = args.version
    if not check_version_format(version):
        print(f'Error: version "{version}" is not in correct format,'
              ' expected "N.N.N" or "vN.N.N"', file=sys.stderr)
        sys.exit(1)

    version = version.lstrip('v')
    remote = args.push
    force = args.force
    tag = f'v{version}'

    if not force and check_tag_exists(tag):
        print(f'Error: tag "{tag}" already exists', file=sys.stderr)
        sys.exit(1)

    write_version(version)
    commit_change(version, force)
    create_tag(tag, force)
    if remote:
        push_change(remote, tag, force)
    print(f'--- Successfully released {tag}')

if __name__ == '__main__':
    main()
