package templates

var CleanupNBDB = `
<<EOF
nbstatus=$(ovs-appctl -t /var/run/ovn/ovnnb_db.ctl cluster/status OVN_Northbound)
echo "current northbound status"
echo "$nbstatus"
nodeID=$(grep '{{ .NodeAddress }}' $nbstatus | awk '{print $1}')
if [ -n "$nodeID" ]
then
  ovs-appctl -t /var/run/ovn/ovnnb_db.ctl cluster/kick OVN_Northbound $nodeID
  echo "removed node id $nodeID with address {{ .NodeAddress }}"
  echo "current northbound status"
  ovs-appctl -t /var/run/ovn/ovnnb_db.ctl cluster/status OVN_Northbound
fi
`

var CleanupSBDB = `
sbStatus=$(ovs-appctl -t /var/run/ovn/ovnsb_db.ctl cluster/status OVN_Southbound)
echo "current southbound status"
echo "$sbStatus"
nodeID=$(grep '{{ .NodeAddress }}' $sbStatus | awk '{print $1}')
if [ -n "$nodeID" ]
then
  ovs-appctl -t /var/run/ovn/ovnsb_db.ctl cluster/kick OVN_Southbound $nodeID
  echo "removed node id $nodeID with address {{ .NodeAddress }}"
  echo "current southbound status"
  ovs-appctl -t /var/run/ovn/ovnsb_db.ctl cluster/status OVN_Southbound
fi
`
