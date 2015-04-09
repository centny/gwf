fvm_s=$1
o_dir=build
rm -rf $o_dir
mkdir $o_dir

#### Package cfgs ####
v_cfgs=0.0.1
o_cfgs=$o_dir/cfgs
go build -o $o_cfgs/cfgs github.com/Centny/gwf/cmd/cfgs
cp cmd/cfgs/cfgsd $o_cfgs
cp cmd/cfgs/token.properties $o_cfgs
cd $o_dir
zip -r cfgs.zip cfgs
cd ../
if [ "$fvm_s" != "" ];then
	fvm -u $fvm_s cfgs $v_cfgs $o_dir/cfgs.zip
fi
echo "Package cfgs..."