#!/usr/bin/env sh

cd $(go env GOROOT)/src/runtime

cat > proc_id.go <<EOF
package runtime

func GoID() int64 {
    return getg().goid
}
EOF

go install
