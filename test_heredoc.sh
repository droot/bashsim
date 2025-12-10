
echo "Testing heredoc..."
./bashsim <<EOF
cat <<END
Hello Heredoc
END
exit
EOF
